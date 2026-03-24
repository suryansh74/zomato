import React, { useState } from "react";
import axios from "axios";
import { restaurantServiceUrl } from "@/lib/config";
import type { Restaurant } from "@/types/types";
import toast from "react-hot-toast";

interface EditRestaurantProps {
  restaurant: Restaurant;
  onCancel: () => void;
  // CHANGED: We now send the updated data back to the parent
  onSaveSuccess: (
    serverData: any,
    localData: { name: string; description: string },
  ) => void;
}

export function EditRestaurant({
  restaurant,
  onCancel,
  onSaveSuccess,
}: EditRestaurantProps) {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: restaurant.name || "",
    description: restaurant.description || "",
  });

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    try {
      setLoading(true);
      const res = await axios.put(
        `${restaurantServiceUrl}/restaurant/update`,
        formData,
        { withCredentials: true },
      );

      toast.success("Restaurant updated successfully");

      // Send both the server response and our local form data back to the parent
      onSaveSuccess(res.data?.data, formData);
    } catch (err) {
      if (axios.isAxiosError(err)) {
        toast.error(
          err.response?.data?.error || "Failed to update restaurant.",
        );
      } else {
        toast.error("An unexpected error occurred.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-8 space-y-5 animate-in slide-in-from-right-2 duration-300">
      <h2 className="text-xl font-bold text-gray-800 mb-4 border-b pb-2">
        Edit Restaurant Details
      </h2>

      <div>
        <label className="block text-sm font-semibold text-gray-700 mb-1">
          Restaurant Name
        </label>
        <input
          type="text"
          name="name"
          required
          value={formData.name}
          onChange={handleInputChange}
          className="w-full text-lg text-gray-800 border border-gray-300 rounded-md px-4 py-2 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-colors"
        />
      </div>

      <div>
        <label className="block text-sm font-semibold text-gray-700 mb-1">
          Description
        </label>
        <textarea
          name="description"
          value={formData.description}
          onChange={handleInputChange}
          rows={4}
          className="w-full text-gray-700 border border-gray-300 rounded-md px-4 py-3 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 resize-none transition-colors"
        />
      </div>

      <div className="flex items-center justify-end gap-3 pt-4 border-t border-gray-100">
        <button
          onClick={onCancel}
          disabled={loading}
          className="px-5 py-2.5 text-gray-600 hover:bg-gray-100 rounded-lg font-medium transition-colors"
        >
          Cancel
        </button>
        <button
          onClick={handleSave}
          disabled={loading || !formData.name.trim()}
          className="flex items-center gap-2 bg-[#2563eb] hover:bg-blue-700 text-white px-6 py-2.5 rounded-lg font-medium transition-colors disabled:opacity-70 shadow-sm"
        >
          {loading ? "Saving..." : "Save Changes"}
        </button>
      </div>
    </div>
  );
}
