import React, { useState } from "react";
import axios from "axios";
import { restaurantServiceUrl } from "@/lib/config";
import { useAuth } from "@/context/useAuth";

interface CreateRestaurantFormProps {
  onSuccess: () => void;
  onCancel: () => void;
}

export function CreateRestaurantForm({
  onSuccess,
  onCancel,
}: CreateRestaurantFormProps) {
  const { location, loadingLocation } = useAuth(); // Tap into your existing context
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    name: "",
    description: "",
    phone: "",
    formatted_address: "",
    latitude: "",
    longitude: "",
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

  // Helper to pre-fill location from context
  const handleUseMyLocation = () => {
    if (location) {
      setFormData((prev) => ({
        ...prev,
        formatted_address: location.formattedAddress,
        latitude: location.latitude.toString(),
        longitude: location.longitude.toString(),
      }));
    } else {
      alert(
        "Location not available. Please ensure location services are enabled.",
      );
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    if (!imageFile) {
      setError("A cover image is required for your restaurant.");
      setLoading(false);
      return;
    }

    const payload = new FormData();
    payload.append("name", formData.name);
    payload.append("description", formData.description);
    payload.append("phone", formData.phone);
    payload.append("formatted_address", formData.formatted_address);
    payload.append("latitude", formData.latitude);
    payload.append("longitude", formData.longitude);
    payload.append("image", imageFile);

    try {
      await axios.post(`${restaurantServiceUrl}/restaurant/create`, payload, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      });
      onSuccess();
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to create restaurant.");
      } else {
        setError("An unexpected error occurred.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto bg-white border border-gray-100 rounded-2xl shadow-2xl overflow-hidden">
      <div className="bg-gradient-to-r from-red-600 to-red-700 p-8 text-white">
        <h2 className="text-3xl font-extrabold mb-2">Create Your Restaurant</h2>
        <p className="text-red-100 font-medium">
          Fill out the details below to set up your business profile.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="p-8 space-y-8">
        {error && (
          <div className="p-4 bg-red-50 border-l-4 border-red-600 text-red-700 rounded-r shadow-sm font-medium">
            {error}
          </div>
        )}

        {/* Basic Info */}
        <div className="space-y-6">
          <h3 className="text-lg font-bold text-gray-900 border-b pb-2">
            Basic Information
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-1">
                Restaurant Name <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="name"
                required
                value={formData.name}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 outline-none transition-all bg-gray-50 focus:bg-white"
                placeholder="e.g. The Spicy Kitchen"
              />
            </div>
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-1">
                Phone Number <span className="text-red-500">*</span>
              </label>
              <input
                type="tel"
                name="phone"
                required
                value={formData.phone}
                onChange={handleInputChange}
                className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 outline-none transition-all bg-gray-50 focus:bg-white"
                placeholder="e.g. +91 98765 43210"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-bold text-gray-700 mb-1">
              Description
            </label>
            <textarea
              name="description"
              rows={3}
              value={formData.description}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 outline-none transition-all bg-gray-50 focus:bg-white resize-none"
              placeholder="Tell customers what makes your food special..."
            />
          </div>
        </div>

        {/* Location Info */}
        <div className="space-y-6 pt-2">
          <div className="flex items-center justify-between border-b pb-2">
            <h3 className="text-lg font-bold text-gray-900">
              Location Details
            </h3>
            <button
              type="button"
              onClick={handleUseMyLocation}
              disabled={loadingLocation}
              className="text-sm font-semibold text-red-600 hover:text-red-800 flex items-center gap-1 bg-red-50 px-3 py-1.5 rounded-lg transition-colors"
            >
              📍 {loadingLocation ? "Locating..." : "Use My Current Location"}
            </button>
          </div>

          <div>
            <label className="block text-sm font-bold text-gray-700 mb-1">
              Full Address <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              name="formatted_address"
              required
              value={formData.formatted_address}
              onChange={handleInputChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-red-500 focus:border-red-500 outline-none transition-all bg-gray-50 focus:bg-white mb-4"
              placeholder="123 Food Street, City, State"
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-xs font-bold text-gray-500 mb-1 uppercase tracking-wider">
                  Latitude
                </label>
                <input
                  type="number"
                  step="any"
                  name="latitude"
                  required
                  value={formData.latitude}
                  onChange={handleInputChange}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 outline-none bg-gray-50 font-mono text-sm"
                  placeholder="e.g. 26.9124"
                />
              </div>
              <div>
                <label className="block text-xs font-bold text-gray-500 mb-1 uppercase tracking-wider">
                  Longitude
                </label>
                <input
                  type="number"
                  step="any"
                  name="longitude"
                  required
                  value={formData.longitude}
                  onChange={handleInputChange}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 outline-none bg-gray-50 font-mono text-sm"
                  placeholder="e.g. 75.7873"
                />
              </div>
            </div>
          </div>
        </div>

        {/* Image Upload */}
        <div className="space-y-6 pt-2">
          <h3 className="text-lg font-bold text-gray-900 border-b pb-2">
            Media
          </h3>
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">
              Cover Image <span className="text-red-500">*</span>
            </label>
            <div className="flex items-center gap-6 p-4 border-2 border-dashed border-gray-300 rounded-xl bg-gray-50">
              {imagePreview ? (
                <img
                  src={imagePreview}
                  alt="Preview"
                  className="w-24 h-24 object-cover rounded-xl shadow-sm border border-gray-200"
                />
              ) : (
                <div className="w-24 h-24 bg-gray-200 rounded-xl flex items-center justify-center text-3xl text-gray-400">
                  📸
                </div>
              )}
              <div className="flex-1">
                <input
                  type="file"
                  accept="image/*"
                  required
                  onChange={handleImageChange}
                  className="block w-full text-sm text-gray-600 file:mr-4 file:py-2.5 file:px-6 file:rounded-full file:border-0 file:text-sm file:font-bold file:bg-red-100 file:text-red-700 hover:file:bg-red-200 cursor-pointer transition-colors"
                />
                <p className="text-xs text-gray-500 mt-2">
                  Recommended: High resolution landscape image (PNG, JPG).
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Form Actions */}
        <div className="pt-8 flex items-center justify-end gap-4 border-t border-gray-200">
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="px-6 py-3 text-gray-700 font-bold hover:bg-gray-100 rounded-xl transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={loading}
            className="px-8 py-3 bg-red-600 hover:bg-red-700 text-white font-bold rounded-xl transition-all shadow-md hover:shadow-lg disabled:opacity-70 disabled:cursor-not-allowed flex items-center gap-2"
          >
            {loading ? (
              <>
                <svg
                  className="animate-spin -ml-1 mr-2 h-5 w-5 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                Creating Profile...
              </>
            ) : (
              "Create Profile"
            )}
          </button>
        </div>
      </form>
    </div>
  );
}
