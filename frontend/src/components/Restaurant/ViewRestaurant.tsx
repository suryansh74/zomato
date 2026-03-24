import React from "react";
import type { Restaurant } from "@/types/types";

interface ViewRestaurantProps {
  restaurant: Restaurant;
  onEditClick: () => void;
  onToggleStatus: () => void;
  isToggling: boolean;
}

export function ViewRestaurant({
  restaurant,
  onEditClick,
  onToggleStatus,
  isToggling,
}: ViewRestaurantProps) {
  return (
    <div>
      <div className="p-8 animate-in fade-in duration-300">
        {/* Header & Edit Button */}
        <div className="flex justify-between items-start mb-4">
          <h1 className="text-3xl font-bold text-gray-900">
            {restaurant.name}
          </h1>
          <button
            onClick={onEditClick}
            className="text-gray-500 hover:text-blue-600 p-2 rounded-full hover:bg-blue-50 transition-colors"
            title="Edit Restaurant Details"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
              />
            </svg>
          </button>
        </div>

        {/* Location */}
        <div className="flex items-start gap-2 text-gray-500 mb-6">
          <svg
            className="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
          <p className="text-sm">
            {restaurant.auto_location?.formatted_address || "Location not set"}
          </p>
        </div>

        {/* Description */}
        <div className="text-gray-700 whitespace-pre-wrap mb-8">
          {restaurant.description || (
            <span className="text-gray-400 italic">
              No description provided.
            </span>
          )}
        </div>

        <hr className="border-gray-200 mb-6" />

        {/* Status & Actions */}
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="font-bold tracking-wide">
            {restaurant.is_open ? (
              <span className="text-green-600 uppercase">● Open</span>
            ) : (
              <span className="text-red-600 uppercase">● Closed</span>
            )}
          </div>

          <button
            onClick={onToggleStatus}
            disabled={isToggling}
            className={`px-6 py-2.5 rounded-lg font-medium text-white transition-colors shadow-sm disabled:opacity-70 ${
              restaurant.is_open
                ? "bg-red-500 hover:bg-red-600"
                : "bg-[#16a34a] hover:bg-green-700"
            }`}
          >
            {isToggling
              ? "Updating..."
              : restaurant.is_open
                ? "Close Restaurant"
                : "Open Restaurant"}
          </button>
        </div>
      </div>
    </div>
  );
}
