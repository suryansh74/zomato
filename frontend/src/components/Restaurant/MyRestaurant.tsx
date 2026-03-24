import type { Restaurant } from "@/types/types";

interface MyRestaurantProps {
  restaurant: Restaurant;
}

export default function MyRestaurant({ restaurant }: MyRestaurantProps) {
  return (
    <div className="max-w-4xl mx-auto bg-white border border-gray-100 rounded-3xl shadow-xl overflow-hidden">
      {/* Header Image */}
      <div className="relative h-80 bg-gray-100 border-b border-gray-100">
        {restaurant.image ? (
          <img
            src={restaurant.image}
            alt={restaurant.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex flex-col items-center justify-center text-gray-400">
            <span className="text-4xl mb-2">📸</span>
            <span>No Image Provided</span>
          </div>
        )}

        {/* Floating Status Badges */}
        <div className="absolute top-4 right-4 flex gap-3">
          <span
            className={`px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-md backdrop-blur-sm ${
              restaurant.is_open
                ? "bg-white/90 text-red-600"
                : "bg-gray-900/90 text-white"
            }`}
          >
            {restaurant.is_open ? "● Open Now" : "Closed"}
          </span>
          {restaurant.is_verified && (
            <span className="px-4 py-1.5 rounded-full text-xs font-bold bg-white/90 text-green-600 shadow-md uppercase tracking-wider backdrop-blur-sm">
              ✓ Verified
            </span>
          )}
        </div>
      </div>

      {/* Content Details */}
      <div className="p-10">
        <h1 className="text-4xl font-black text-gray-900 mb-4 tracking-tight">
          {restaurant.name}
        </h1>

        {restaurant.description && (
          <p className="text-gray-600 text-lg mb-10 leading-relaxed max-w-3xl">
            {restaurant.description}
          </p>
        )}

        {/* Contact & Location Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 pt-8 border-t border-gray-100 bg-gray-50/50 rounded-2xl p-6">
          <div className="flex items-start gap-4">
            <div className="p-4 bg-red-50 text-red-600 rounded-xl text-2xl shadow-sm">
              📞
            </div>
            <div>
              <p className="text-sm font-bold text-gray-500 mb-1 uppercase tracking-wider">
                Phone Number
              </p>
              <p className="text-gray-900 font-semibold text-lg">
                {restaurant.phone}
              </p>
            </div>
          </div>

          <div className="flex items-start gap-4">
            <div className="p-4 bg-red-50 text-red-600 rounded-xl text-2xl shadow-sm">
              📍
            </div>
            <div>
              <p className="text-sm font-bold text-gray-500 mb-1 uppercase tracking-wider">
                Location
              </p>
              <p className="text-gray-900 font-medium leading-snug">
                {restaurant.auto_location?.formatted_address}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
