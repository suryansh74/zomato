import React, { useState } from "react";
import axios from "axios";
import { restaurantServiceUrl } from "@/lib/config";
import { useAuth } from "@/context/useAuth";
import toast from "react-hot-toast";

interface CreateRestaurantFormProps {
  onSuccess: () => void;
  onCancel: () => void; // Keeping this if you still want a cancel option, otherwise can be removed
}

export function CreateRestaurantForm({
  onSuccess,
  onCancel,
}: CreateRestaurantFormProps) {
  const { location, loadingLocation } = useAuth(); // Automatically fetched from AuthProvider
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    name: "",
    phone: "",
    description: "",
  });

  const [imageFile, setImageFile] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setImageFile(file);
      setImagePreview(URL.createObjectURL(file));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // 1. Validations
    if (!imageFile) {
      setError("Please upload a restaurant image.");
      setLoading(false);
      return;
    }

    if (!location) {
      setError(
        "Location is required. Please allow location access in your browser.",
      );
      setLoading(false);
      return;
    }

    // 2. Prepare Payload
    const payload = new FormData();
    payload.append("name", formData.name);
    payload.append("phone", formData.phone);
    payload.append("description", formData.description);

    // Append location data securely from context (user cannot edit this)
    payload.append("formatted_address", location.formattedAddress);
    payload.append("latitude", location.latitude.toString());
    payload.append("longitude", location.longitude.toString());

    // Append the file
    payload.append("image", imageFile);

    // 3. Submit
    try {
      await axios.post(`${restaurantServiceUrl}/restaurant/create`, payload, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      });
      toast.success("Restaurant added successfully!");
      onSuccess();
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to add restaurant.");
      } else {
        setError("An unexpected error occurred.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-xl mx-auto bg-white p-8 rounded-xl shadow-sm border border-gray-100 mt-10">
      <h2 className="text-2xl font-semibold text-gray-900 mb-8">
        Add Your Restaurant
      </h2>

      <form onSubmit={handleSubmit} className="space-y-5">
        {error && (
          <div className="p-3 bg-red-50 border-l-4 border-red-500 text-red-700 text-sm rounded">
            {error}
          </div>
        )}

        {/* Name */}
        <div>
          <input
            type="text"
            name="name"
            required
            value={formData.name}
            onChange={handleInputChange}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-1 focus:ring-red-500 focus:border-red-500 outline-none placeholder-gray-400"
            placeholder="Restaurant name"
          />
        </div>

        {/* Contact Number */}
        <div>
          <input
            type="tel"
            name="phone"
            required
            value={formData.phone}
            onChange={handleInputChange}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-1 focus:ring-red-500 focus:border-red-500 outline-none placeholder-gray-400"
            placeholder="Contact Number"
          />
        </div>

        {/* Description */}
        <div>
          <textarea
            name="description"
            rows={4}
            value={formData.description}
            onChange={handleInputChange}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-1 focus:ring-red-500 focus:border-red-500 outline-none placeholder-gray-400 resize-none"
            placeholder="Restaurant Description"
          />
        </div>

        {/* Image Upload Box */}
        <div className="relative border border-gray-300 rounded-lg p-4 flex items-center gap-3 hover:bg-gray-50 transition-colors cursor-pointer group overflow-hidden">
          <input
            type="file"
            accept="image/*"
            required
            onChange={handleImageChange}
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
          />
          {imagePreview ? (
            <div className="flex items-center gap-4 w-full">
              <img
                src={imagePreview}
                alt="Preview"
                className="w-12 h-12 object-cover rounded-md"
              />
              <span className="text-gray-600 font-medium text-sm truncate">
                {imageFile?.name}
              </span>
            </div>
          ) : (
            <>
              <svg
                className="w-5 h-5 text-red-500"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                />
              </svg>
              <span className="text-gray-500 group-hover:text-gray-700">
                Upload restaurant image
              </span>
            </>
          )}
        </div>

        {/* Auto Fetched Read-Only Location */}
        <div className="flex items-start gap-3 text-gray-700 py-4">
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
          <p className="text-sm leading-relaxed">
            {loadingLocation ? (
              <span className="animate-pulse text-gray-400">
                Locating your restaurant...
              </span>
            ) : location ? (
              location.formattedAddress
            ) : (
              <span className="text-red-500">
                Location not available. Please allow access.
              </span>
            )}
          </p>
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          disabled={loading || loadingLocation || !location}
          className="w-full py-3.5 bg-[#df3b4c] hover:bg-red-700 text-white font-medium rounded-lg transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed mt-2"
        >
          {loading ? "Adding..." : "Add Restaurant"}
        </button>

        {/* Optional Cancel Button (If needed based on parent structure) */}
        <button
          type="button"
          onClick={onCancel}
          className="w-full py-2 text-gray-500 hover:text-gray-700 font-medium text-sm text-center"
        >
          Cancel
        </button>
      </form>
    </div>
  );
}
