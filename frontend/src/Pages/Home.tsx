import React, { useState, useEffect } from "react";
import { useSearchParams, useNavigate, Link } from "react-router-dom";
import axios from "axios";
import toast from "react-hot-toast";
import { MapPin, Store, Loader2 } from "lucide-react"; // Removed Search icon

import { restaurantServiceUrl } from "@/lib/config";
import type { Restaurant } from "@/types/types";
import { useAuth } from "@/context/useAuth"; // Assuming your context path

import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
// Removed Input import since we don't need the search bar here anymore

// Helper function to format distance beautifully
const formatDistance = (distanceKm?: number) => {
  if (distanceKm === undefined) return null;
  if (distanceKm < 1) {
    return `${Math.round(distanceKm * 1000)} m`;
  }
  return `${distanceKm.toFixed(1)} km`;
};

export default function Home() {
  const { location, loadingLocation } = useAuth();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  // This automatically reads what your Navbar types into the URL!
  const search = searchParams.get("search") || "";

  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [loading, setLoading] = useState(true);

  // Handle Fresh Login Toast
  useEffect(() => {
    if (searchParams.get("fresh") === "true") {
      toast.success("Logged in successfully!", { id: "login-success" });
      navigate("/", { replace: true });
    }
  }, [searchParams, navigate]);

  // Fetch Nearby Restaurants
  const fetchRestaurants = async () => {
    if (!location?.latitude || !location?.longitude) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);

      const urlRadius = searchParams.get("radius");
      const radiusToUse = urlRadius ? Number(urlRadius) : 10000;

      const response = await axios.get(
        `${restaurantServiceUrl}/restaurant/nearby`,
        {
          params: {
            latitude: location.latitude,
            longitude: location.longitude,
            search: search, // Passes the search term from the Navbar to the Go backend
            radius: radiusToUse,
          },
          withCredentials: true,
        },
      );

      // Safely unwrap Axios data -> Go response data -> restaurants array
      const fetchedRestaurants = response.data?.data?.restaurants || [];
      setRestaurants(fetchedRestaurants);
    } catch (error) {
      console.error("Failed to fetch nearby restaurants:", error);
      toast.error("Failed to load nearby restaurants.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!loadingLocation) {
      fetchRestaurants();
    }
  }, [location, search, loadingLocation]); // Re-runs anytime the Navbar search changes!

  return (
    <div className="max-w-6xl mx-auto px-4 py-8 animate-in fade-in duration-500">
      {/* Header (Search bar removed!) */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">
          Food delivery in your area
        </h1>
      </div>

      {/* STATE 1: Waiting for user to allow location */}
      {loadingLocation && (
        <div className="p-6 bg-blue-50 border border-blue-200 rounded-xl text-blue-800 flex items-center gap-3">
          <Loader2 className="w-6 h-6 animate-spin" />
          <p className="font-medium">Acquiring your location...</p>
        </div>
      )}

      {/* STATE 2: Location denied or unavailable */}
      {!loadingLocation && !location?.latitude && (
        <div className="p-6 bg-amber-50 border border-amber-200 rounded-xl text-amber-800 flex items-center gap-3">
          <MapPin className="w-6 h-6" />
          <p className="font-medium">
            Please allow location access in your browser to see restaurants near
            you.
          </p>
        </div>
      )}

      {/* STATE 3: Loading Restaurants Skeleton */}
      {loading && !loadingLocation && location?.latitude && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map((n) => (
            <div key={n} className="space-y-3">
              <Skeleton className="h-48 w-full rounded-2xl" />
              <Skeleton className="h-5 w-3/4" />
              <Skeleton className="h-4 w-1/2" />
            </div>
          ))}
        </div>
      )}

      {/* STATE 4: Empty Results */}
      {!loading &&
        !loadingLocation &&
        location?.latitude &&
        restaurants.length === 0 && (
          <div className="text-center py-20">
            <Store className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-gray-900 mb-2">
              No restaurants found
            </h3>
            <p className="text-gray-500">
              Try adjusting your search or increasing your radius.
            </p>
          </div>
        )}

      {/* STATE 5: Restaurant Grid */}
      {!loading && restaurants.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {restaurants.map((restaurant) => (
            <Link
              to={`/restaurant/${restaurant.id}`}
              key={restaurant.id}
              className="block group"
            >
              <Card className="border-0 shadow-none overflow-hidden bg-transparent group-hover:bg-gray-50 transition-colors rounded-2xl">
                <CardContent className="p-0">
                  {/* Image Container */}
                  <div className="relative h-48 w-full rounded-2xl overflow-hidden bg-gray-100 mb-3">
                    {restaurant.image ? (
                      <img
                        src={restaurant.image}
                        alt={restaurant.name}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-gray-400">
                        <Store className="w-10 h-10" />
                      </div>
                    )}

                    {!restaurant.is_open && (
                      <div className="absolute inset-0 bg-black/40 flex items-center justify-center backdrop-blur-[2px]">
                        <span className="bg-white text-gray-900 text-xs font-bold uppercase tracking-wider px-3 py-1.5 rounded-full shadow-md">
                          Currently Closed
                        </span>
                      </div>
                    )}
                  </div>

                  {/* Details */}
                  <div>
                    <div className="flex justify-between items-start">
                      <h3
                        className={`text-lg font-bold line-clamp-1 ${!restaurant.is_open ? "text-gray-500" : "text-gray-900"}`}
                      >
                        {restaurant.name}
                      </h3>
                      {/* Distance formatting applied here! */}
                      {restaurant.distance_km !== undefined && (
                        <span className="text-sm font-semibold text-gray-600 bg-gray-100 px-2 py-0.5 rounded-md whitespace-nowrap ml-2">
                          {formatDistance(restaurant.distance_km)}
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-500 line-clamp-1 mt-0.5">
                      {restaurant.auto_location?.formatted_address ||
                        "Location unavailable"}
                    </p>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
