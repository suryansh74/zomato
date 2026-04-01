// Restaurant.tsx — Page component
import { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { Link } from "react-router-dom"; // ✅ Added Link
import { MonitorPlay } from "lucide-react"; // ✅ Added Icon
import { restaurantServiceUrl } from "@/lib/config";
import type { Restaurant as RestaurantType } from "@/types/types";
import MyRestaurant from "@/components/Restaurant/MyRestaurant";
import { CreateRestaurantForm } from "@/components/Restaurant/CreateRestaurantForm";

export default function Restaurant() {
  const [restaurant, setRestaurant] = useState<RestaurantType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);

  const fetchMyRestaurant = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      // ✅ Using YOUR correct endpoint to prevent the 404
      const res = await axios.get(`${restaurantServiceUrl}/restaurant/read`, {
        withCredentials: true,
      });
      const data: RestaurantType = res.data?.data ?? res.data;
      setRestaurant(data);
    } catch (err) {
      if (axios.isAxiosError(err)) {
        if (err.response?.status === 404) {
          setRestaurant(null);
        } else {
          setError(err.response?.data?.error ?? "Failed to fetch restaurant");
        }
      } else {
        setError("An unexpected error occurred.");
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMyRestaurant();
  }, [fetchMyRestaurant]);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50">
        <div className="animate-pulse text-red-600 font-semibold text-xl">
          Loading dashboard...
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-3xl mx-auto mt-12 p-6 bg-red-50 border-l-4 border-red-600 text-red-800 rounded-lg shadow-sm">
        <h3 className="font-bold text-lg mb-2">Error Loading Dashboard</h3>
        <p>{error}</p>
        <button
          onClick={fetchMyRestaurant}
          className="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
        >
          Try Again
        </button>
      </div>
    );
  }

  if (!restaurant) {
    if (isCreating) {
      return (
        <div className="min-h-screen bg-gray-50 py-12 px-4">
          <CreateRestaurantForm
            onSuccess={() => {
              setIsCreating(false);
              fetchMyRestaurant();
            }}
            onCancel={() => setIsCreating(false)}
          />
        </div>
      );
    }
    return (
      <div className="min-h-screen bg-gray-50 py-20 px-4">
        <div className="max-w-2xl mx-auto bg-white border border-gray-200 rounded-2xl shadow-sm p-12 text-center">
          <div className="w-20 h-20 bg-red-50 text-red-600 rounded-full flex items-center justify-center mx-auto mb-6 text-4xl shadow-inner">
            🍽️
          </div>
          <h2 className="text-3xl font-extrabold text-gray-900 mb-4">
            Welcome to Your Dashboard
          </h2>
          <p className="text-gray-500 text-lg mb-10 max-w-lg mx-auto">
            You haven't set up a restaurant profile yet. Create one now to start
            managing your menu, orders, and business details.
          </p>
          <button
            onClick={() => setIsCreating(true)}
            className="px-8 py-4 bg-red-600 hover:bg-red-700 text-white text-lg font-semibold rounded-xl transition-all shadow-md hover:shadow-lg hover:-translate-y-0.5"
          >
            + Create Restaurant Profile
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      {/* ✅ NEW: Dashboard Button injected right above your existing UI */}
      <div className="max-w-4xl mx-auto mb-6 flex justify-end">
        <Link
          to="/restaurant/dashboard"
          className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-xl font-bold transition shadow-md animate-in fade-in slide-in-from-top-4"
        >
          <MonitorPlay className="w-5 h-5" />
          Open Live Kitchen Dashboard
        </Link>
      </div>

      {/* Your existing MyRestaurant component */}
      <MyRestaurant restaurant={restaurant} />
    </div>
  );
}
