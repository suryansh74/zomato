import React, { useState, useEffect } from "react";
import axios from "axios";
import { restaurantServiceUrl } from "@/lib/config";
import type { Restaurant } from "@/types/types";
import toast from "react-hot-toast";

import { ViewRestaurant } from "./ViewRestaurant";
import { EditRestaurant } from "./EditRestaurant";

interface MyRestaurantProps {
  restaurant: Restaurant;
  onRefresh?: () => void;
}

export default function MyRestaurant({
  restaurant: initialRestaurant,
  onRefresh,
}: MyRestaurantProps) {
  const [currentRestaurant, setCurrentRestaurant] =
    useState<Restaurant>(initialRestaurant);
  const [isEditing, setIsEditing] = useState(false);
  const [isToggling, setIsToggling] = useState(false);

  // Sync state if parent actually passes down new props
  useEffect(() => {
    setCurrentRestaurant(initialRestaurant);
  }, [initialRestaurant]);

  // Handle immediate Open/Close toggle
  const handleToggleStatus = async () => {
    const newStatus = !currentRestaurant.is_open;
    try {
      setIsToggling(true);
      const res = await axios.put(
        `${restaurantServiceUrl}/restaurant/update`,
        { is_open: newStatus },
        { withCredentials: true },
      );

      // INSTANT UI UPDATE: Safely update local state immediately
      const serverData = res.data?.data;
      if (serverData && serverData.id) {
        setCurrentRestaurant(serverData);
      } else {
        setCurrentRestaurant((prev) => ({ ...prev, is_open: newStatus }));
      }

      toast.success(
        newStatus ? "Restaurant is now Open!" : "Restaurant is now Closed.",
      );

      // Still refresh in background to keep parent in sync
      if (onRefresh) onRefresh();
    } catch (err) {
      console.error(err);
      toast.error("Failed to change restaurant status.");
    } finally {
      setIsToggling(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden mt-8 mb-12">
      <div className="h-72 w-full bg-gray-100 relative border-b border-gray-200">
        {currentRestaurant.image ? (
          <img
            src={currentRestaurant.image}
            alt={currentRestaurant.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="flex items-center justify-center h-full text-gray-400">
            No Image Provided
          </div>
        )}
      </div>

      {!isEditing ? (
        <ViewRestaurant
          restaurant={currentRestaurant}
          onEditClick={() => setIsEditing(true)}
          onToggleStatus={handleToggleStatus}
          isToggling={isToggling}
        />
      ) : (
        <EditRestaurant
          restaurant={currentRestaurant}
          onCancel={() => setIsEditing(false)}
          onSaveSuccess={(serverData, localFormData) => {
            // INSTANT UI UPDATE: We update the state before even closing the form
            if (serverData && serverData.id) {
              setCurrentRestaurant(serverData); // Use backend data if perfect
            } else {
              setCurrentRestaurant((prev) => ({ ...prev, ...localFormData })); // Fallback to what we just typed
            }

            setIsEditing(false); // Switch back to view mode

            // Background refresh
            if (onRefresh) onRefresh();
          }}
        />
      )}
    </div>
  );
}
