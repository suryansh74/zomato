import React, { useState, useEffect } from "react";
import axios from "axios";
import { restaurantServiceUrl } from "@/lib/config";
import type { Restaurant } from "@/types/types";
import toast from "react-hot-toast";

import { ViewRestaurant } from "./ViewRestaurant";
import { EditRestaurant } from "./EditRestaurant";
import { MenuDashboard } from "./MenuDashboard"; // Import it here instead!

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

  useEffect(() => {
    setCurrentRestaurant(initialRestaurant);
  }, [initialRestaurant]);

  const handleToggleStatus = async () => {
    const newStatus = !currentRestaurant.is_open;
    try {
      setIsToggling(true);
      const res = await axios.put(
        `${restaurantServiceUrl}/restaurant/update`,
        { is_open: newStatus },
        { withCredentials: true },
      );

      const serverData = res.data?.data;
      if (serverData && serverData.id) {
        setCurrentRestaurant(serverData);
      } else {
        setCurrentRestaurant((prev) => ({ ...prev, is_open: newStatus }));
      }

      toast.success(
        newStatus ? "Restaurant is now Open!" : "Restaurant is now Closed.",
      );
      if (onRefresh) onRefresh();
    } catch (err) {
      toast.error("Failed to change restaurant status.");
    } finally {
      setIsToggling(false);
    }
  };

  return (
    // Note: I removed the max-w-3xl from this outermost wrapper so the layout has more breathing room,
    // but kept the restaurant card itself constrained.
    <div className="w-full max-w-5xl mx-auto px-4 pb-16">
      {/* 1. THE RESTAURANT DETAILS CARD */}
      <div className="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden mt-8 mb-12 max-w-3xl mx-auto">
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
              if (serverData && serverData.id) {
                setCurrentRestaurant(serverData);
              } else {
                setCurrentRestaurant((prev) => ({ ...prev, ...localFormData }));
              }
              setIsEditing(false);
              if (onRefresh) onRefresh();
            }}
          />
        )}
      </div>

      {/* 2. THE MENU DASHBOARD */}
      {/* Placed OUTSIDE the restaurant card so the tabs can stretch nicely */}
      <div className="max-w-4xl mx-auto">
        <MenuDashboard />
      </div>
    </div>
  );
}
