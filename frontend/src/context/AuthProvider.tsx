import { useEffect, useState } from "react";
import type { LocationData, User } from "../types/types";
import { AuthContext } from "./AuthContext";

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [location, setLocation] = useState<LocationData | null>(null);
  const [loadingLocation, setLoadingLocation] = useState(false);
  const [city, setCity] = useState("Fetching Location...");

  useEffect(() => {
    fetch("http://localhost:8000/api/auth/profile", {
      credentials: "include",
    })
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => setUser(data?.data?.payload?.user ?? null))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!navigator.geolocation) {
      alert("Please allow location to continue");
      return;
    }

    // don't call setState directly here — use a flag instead
    navigator.geolocation.getCurrentPosition(
      async (position) => {
        setLoadingLocation(true); // ✅ inside callback, not effect body
        const { latitude, longitude } = position.coords;
        try {
          const res = await fetch(
            `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${latitude}&lon=${longitude}`,
          );
          const data = await res.json();
          setLocation({
            latitude,
            longitude,
            formattedAddress: data.display_name || "Custom Location",
          });
          setCity(
            data.address.city ||
              data.address.town ||
              data.address.village ||
              "Your Location",
          );
          setLoadingLocation(false);
        } catch (error) {
          setLocation({
            latitude,
            longitude,
            formattedAddress: "Custom Location",
          });
          setCity("Failed to Load");
          setLoadingLocation(false);
          console.error(error);
        } finally {
          setLoadingLocation(false);
        }
      },
      // error callback for when user denies location
      () => {
        setCity("Location denied");
        setLoadingLocation(false);
      },
    );
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, setUser, loading, location, city, loadingLocation }}
    >
      {children}
    </AuthContext.Provider>
  );
}
