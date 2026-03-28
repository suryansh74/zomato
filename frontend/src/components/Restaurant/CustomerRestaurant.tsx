import { useState, useEffect, useCallback } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import toast from "react-hot-toast";
import { MapPin, Plus } from "lucide-react";

import { restaurantServiceUrl } from "@/lib/config";
import type { Restaurant, MenuItem } from "@/types/types";
import { useCart } from "@/context/useCart";

export default function CustomerRestaurant() {
  const { id } = useParams<{ id: string }>(); // Grabs the ID from the URL
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);

  const { setCartLength, fetchCart } = useCart(); // <-- PULL IN SETTER
  const fetchData = useCallback(async () => {
    try {
      setLoading(true);

      // 1. Fetch the Restaurant details
      const resData = await axios.get(
        `${restaurantServiceUrl}/restaurant/${id}`,
        { withCredentials: true },
      );

      console.log("Raw Restaurant Response:", resData.data); // <-- Look at this in your console!

      // THE FIX: Unwrap Axios (data) -> Unwrap Go (data) -> Grab restaurant
      const fetchedRestaurant = resData.data?.data?.restaurant;
      setRestaurant(fetchedRestaurant || null);

      // 2. Fetch the Menu Items
      const menuRes = await axios.get(
        `${restaurantServiceUrl}/restaurant/${id}/menu`,
        { withCredentials: true },
      );

      console.log("Raw Menu Response:", menuRes.data); // <-- Look at this in your console!

      // THE FIX: Unwrap Axios (data) -> Unwrap Go (data) -> Grab menu_items
      const fetchedMenu = menuRes.data?.data?.menu_items || [];
      setMenuItems(fetchedMenu);
    } catch (error) {
      console.error("Failed to fetch data:", error);
      toast.error("Failed to load restaurant details.");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    if (id) fetchData();
  }, [fetchData, id]);

  const handleAddToCart = async (item: MenuItem) => {
    if (!restaurant?.id) {
      toast.error("Restaurant information is missing.");
      return;
    }

    try {
      const response = await axios.post(
        `${restaurantServiceUrl}/add_to_cart`,
        {
          restaurantId: restaurant.id,
          itemId: item.id,
        },
        {
          withCredentials: true,
        },
      );

      toast.success(response.data.message || `Added ${item.name} to cart!`);

      // ✅ THE FIX: Tell the global CartProvider to pull the fresh data from the backend!
      // This automatically updates the navbar number AND the cart items array simultaneously.
      fetchCart();
    } catch (error: any) {
      console.error("Add to cart error:", error);

      if (error.response?.status === 409) {
        toast.error(
          "You can order from only one restaurant at a time. Please clear your cart first.",
          { duration: 5000 },
        );
      } else {
        toast.error(
          error.response?.data?.error || "Failed to add item to cart.",
        );
      }
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        Loading...
      </div>
    );
  }

  if (!restaurant) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        Restaurant not found.
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* 1. RESTAURANT HERO SECTION (Read-Only) */}
      <div className="bg-white border-b">
        <div className="max-w-6xl mx-auto px-4 py-8">
          <div className="flex flex-col md:flex-row gap-8 items-start">
            {/* Image */}
            <div className="w-full md:w-1/3 h-64 bg-gray-100 rounded-2xl overflow-hidden shrink-0">
              {restaurant.image ? (
                <img
                  src={restaurant.image}
                  alt={restaurant.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-gray-400">
                  No Image
                </div>
              )}
            </div>

            {/* Info */}
            <div className="flex-1 space-y-4">
              <div className="flex justify-between items-start">
                <h1 className="text-4xl font-bold text-gray-900">
                  {restaurant.name}
                </h1>
                <span
                  className={`px-4 py-1.5 rounded-full text-sm font-bold uppercase tracking-wide ${restaurant.is_open ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"}`}
                >
                  {restaurant.is_open ? "Open Now" : "Closed"}
                </span>
              </div>

              <div className="flex items-center gap-2 text-gray-500">
                <MapPin className="w-5 h-5 text-red-500" />
                <p>
                  {restaurant.auto_location?.formatted_address ||
                    "Location unavailable"}
                </p>
              </div>

              <div className="text-gray-700">
                {restaurant.description || "No description provided."}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 2. MENU SECTION */}
      <div className="max-w-6xl mx-auto px-4 py-12">
        <h2 className="text-2xl font-bold text-gray-900 mb-8">Order Online</h2>

        {menuItems.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
            <p className="text-gray-500">
              This restaurant hasn't added any menu items yet.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {menuItems.map((item) => (
              <div
                key={item.id}
                className="bg-white border rounded-xl overflow-hidden hover:shadow-md transition-shadow p-4 flex gap-4"
              >
                {/* Item Details */}
                <div className="flex-1">
                  <h3 className="font-bold text-gray-900 text-lg mb-1">
                    {item.name}
                  </h3>
                  <p className="font-semibold text-gray-700 mb-2">
                    ₹{item.price}
                  </p>
                  <p className="text-sm text-gray-500 line-clamp-2">
                    {item.description}
                  </p>
                </div>

                {/* Image & Add Button */}
                <div className="w-32 flex flex-col items-center gap-3 shrink-0">
                  <div className="w-full h-32 bg-gray-100 rounded-lg overflow-hidden">
                    {item.image ? (
                      <img
                        src={item.image}
                        alt={item.name}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-xs text-gray-400">
                        No Image
                      </div>
                    )}
                  </div>

                  <button
                    onClick={() => handleAddToCart(item)}
                    disabled={!item.is_available || !restaurant.is_open}
                    className="w-full py-2 bg-red-50 text-red-600 font-bold rounded-lg hover:bg-red-100 transition-colors flex items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <Plus className="w-4 h-4" /> ADD
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
