import {
  MapContainer,
  TileLayer,
  Marker,
  useMapEvents,
  useMap,
} from "react-leaflet";
import { useEffect, useState } from "react";
import axios from "axios";
import toast from "react-hot-toast";
import L from "leaflet";
import { Crosshair, Loader2, Plus, Trash2, MapPin, Phone } from "lucide-react"; // ✅ Switched to Lucide

import { restaurantServiceUrl } from "@/lib/config"; // ✅ Using your config

// 🔧 Fix leaflet marker icon issue
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
  iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
  shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
});

// ✅ Matched the Go Backend Struct
interface Address {
  id: string;
  formatted_address: string;
  mobile: string;
}

// 📍 Click-to-select location
const LocationPicker = ({
  setLocation,
}: {
  setLocation: (lat: number, lng: number) => void;
}) => {
  useMapEvents({
    click(e) {
      setLocation(e.latlng.lat, e.latlng.lng);
    },
  });
  return null;
};

// 🎯 Locate me button
const LocateMeButton = ({
  onLocate,
}: {
  onLocate: (lat: number, lng: number) => void;
}) => {
  const map = useMap();
  const locateUser = () => {
    if (!navigator.geolocation) {
      toast.error("Geolocation not supported");
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        map.flyTo([latitude, longitude], 16, { animate: true });
        onLocate(latitude, longitude);
      },
      () => toast.error("Location permission denied"),
    );
  };
  return (
    <button
      onClick={locateUser}
      className="absolute right-3 top-3 z-[1000] flex items-center gap-2 rounded-lg bg-white px-3 py-2 text-sm shadow hover:bg-gray-100"
    >
      <Crosshair size={16} className="text-red-500" />
      Use current location
    </button>
  );
};

export default function AddAddressPage() {
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  // 📋 Form state
  const [mobile, setMobile] = useState("");
  const [formattedAddress, setFormattedAddress] = useState("");
  const [latitude, setLatitude] = useState<number | null>(null);
  const [longitude, setLongitude] = useState<number | null>(null);

  // 🌍 Reverse geocoding
  const fetchFormattedAddress = async (lat: number, lng: number) => {
    try {
      const res = await fetch(
        `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`,
      );
      const data = await res.json();
      setFormattedAddress(data.display_name || "");
    } catch {
      toast.error("Failed to fetch address");
    }
  };

  const setLocation = (lat: number, lng: number) => {
    setLatitude(lat);
    setLongitude(lng);
    fetchFormattedAddress(lat, lng);
  };

  // 📡 Fetch addresses
  const fetchAddresses = async () => {
    try {
      const res = await axios.get(`${restaurantServiceUrl}/address`, {
        withCredentials: true,
      });

      // ✅ THE FIX: Unwrap both the Axios 'data' and the Go 'data' wrappers!
      const fetchedAddresses =
        res.data?.data?.addresses || res.data?.addresses || [];

      setAddresses(fetchedAddresses);
    } catch {
      toast.error("Failed to load addresses");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAddresses();
  }, []);

  // ➕ Add address
  const addAddress = async () => {
    if (
      !mobile ||
      !formattedAddress ||
      latitude === null ||
      longitude === null
    ) {
      toast.error("Please select location on map and enter mobile number");
      return;
    }
    try {
      setAdding(true);
      // ✅ Send exactly what the Go AddressRequest struct expects
      await axios.post(
        `${restaurantServiceUrl}/address`,
        {
          formatted_address: formattedAddress,
          mobile,
          latitude,
          longitude,
        },
        {
          withCredentials: true,
        },
      );
      toast.success("Address added");
      setMobile("");
      setFormattedAddress("");
      setLatitude(null);
      setLongitude(null);
      fetchAddresses();
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to add address");
    } finally {
      setAdding(false);
    }
  };

  // 🗑 Delete address
  const deleteAddress = async (id: string) => {
    if (!window.confirm("Delete this address?")) return;
    try {
      setDeletingId(id);
      await axios.delete(`${restaurantServiceUrl}/address/${id}`, {
        withCredentials: true,
      });
      toast.success("Address deleted");
      fetchAddresses();
    } catch {
      toast.error("Failed to delete address");
    } finally {
      setDeletingId(null);
    }
  };

  // Ensure leaflet CSS is loaded
  useEffect(() => {
    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = "https://unpkg.com/leaflet@1.9.4/dist/leaflet.css";
    document.head.appendChild(link);
  }, []);

  return (
    <div className="mx-auto max-w-4xl px-4 py-8 space-y-8">
      <h1 className="text-3xl font-bold text-gray-900">Delivery Address</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* LEFT COLUMN: Map & Form */}
        <div className="space-y-6">
          <div className="bg-white p-4 rounded-2xl shadow-sm border border-gray-100 space-y-4">
            <h2 className="font-semibold text-gray-800">Add New Address</h2>

            {/* 🗺 Map (Fixed height) */}
            <div className="relative h-[300px] w-full overflow-hidden rounded-xl border border-gray-200 z-0">
              <MapContainer
                center={[latitude || 26.9124, longitude || 75.7873]} // Default Jaipur
                zoom={13}
                className="h-full w-full"
                style={{ height: "100%", width: "100%" }}
              >
                <TileLayer
                  url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                  attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                />
                <LocationPicker setLocation={setLocation} />
                <LocateMeButton onLocate={setLocation} />
                {latitude && longitude && (
                  <Marker position={[latitude, longitude]} />
                )}
              </MapContainer>
            </div>

            {/* 📍 Selected address */}
            {formattedAddress && (
              <div className="flex items-start gap-2 rounded-lg border border-red-100 bg-red-50 p-3 text-sm text-red-900">
                <MapPin className="w-5 h-5 shrink-0 text-red-500 mt-0.5" />
                <p>{formattedAddress}</p>
              </div>
            )}

            {/* 📱 Mobile Input */}
            <div className="relative">
              <Phone className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
              <input
                type="tel" // Go accepts strings, so tel is better than number
                placeholder="Mobile number"
                value={mobile}
                onChange={(e) => setMobile(e.target.value)}
                className="w-full rounded-xl border border-gray-200 pl-10 pr-4 py-3 outline-none focus:border-red-500 focus:ring-1 focus:ring-red-500 transition"
              />
            </div>

            {/* ➕ Save */}
            <button
              disabled={adding || !mobile || !formattedAddress}
              onClick={addAddress}
              className="w-full flex items-center justify-center gap-2 rounded-xl bg-red-500 px-4 py-3 font-bold text-white hover:bg-red-600 disabled:opacity-50 transition"
            >
              {adding ? (
                <Loader2 className="animate-spin w-5 h-5" />
              ) : (
                <Plus className="w-5 h-5" />
              )}
              Save Address
            </button>
          </div>
        </div>

        {/* RIGHT COLUMN: Saved Addresses */}
        <div className="space-y-4">
          <h2 className="text-xl font-bold text-gray-900">Saved Addresses</h2>

          {loading ? (
            <div className="flex items-center gap-2 text-gray-500">
              <Loader2 className="w-5 h-5 animate-spin" /> Loading...
            </div>
          ) : addresses.length === 0 ? (
            <div className="rounded-2xl border border-dashed border-gray-300 p-8 text-center bg-gray-50">
              <MapPin className="w-12 h-12 text-gray-300 mx-auto mb-3" />
              <p className="text-gray-500 font-medium">
                No addresses saved yet
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {addresses.map((addr) => (
                <div
                  key={addr.id} // ✅ Using id instead of _id
                  className="flex items-start justify-between rounded-xl border border-gray-100 bg-white p-4 shadow-sm hover:shadow-md transition"
                >
                  <div className="flex-1 pr-4">
                    <p className="text-sm font-medium text-gray-900 leading-snug">
                      {addr.formatted_address}
                    </p>
                    <p className="text-xs font-semibold text-gray-500 mt-2 flex items-center gap-1">
                      <Phone className="w-3 h-3" /> {addr.mobile}
                    </p>
                  </div>
                  <button
                    onClick={() => deleteAddress(addr.id)}
                    disabled={deletingId === addr.id}
                    className="rounded-lg p-2 text-red-400 hover:bg-red-50 hover:text-red-600 disabled:opacity-50 transition shrink-0"
                  >
                    {deletingId === addr.id ? (
                      <Loader2 size={18} className="animate-spin" />
                    ) : (
                      <Trash2 size={18} />
                    )}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
