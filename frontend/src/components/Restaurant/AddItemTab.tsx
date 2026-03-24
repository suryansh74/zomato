import React, { useState } from "react";
import axios from "axios";
import toast from "react-hot-toast";
import { restaurantServiceUrl } from "@/lib/config";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";

export function AddItemTab({ onSuccess }: { onSuccess: () => void }) {
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    price: "",
    is_available: true,
  });
  const [imageFile, setImageFile] = useState<File | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    const payload = new FormData();
    payload.append("name", formData.name);
    payload.append("description", formData.description);
    payload.append("price", formData.price);
    payload.append("is_available", String(formData.is_available));
    if (imageFile) payload.append("image", imageFile);

    try {
      await axios.post(`${restaurantServiceUrl}/menu`, payload, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      });
      toast.success("Menu item added successfully!");
      setFormData({ name: "", description: "", price: "", is_available: true });
      setImageFile(null);
      onSuccess();
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Failed to add item");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden animate-in fade-in duration-300">
      <div className="p-6 border-b border-gray-100 bg-gray-50/50">
        <h3 className="text-xl font-bold text-gray-900">Add Menu Item</h3>
        <p className="text-sm text-gray-500">
          Create a new dish to display on your restaurant page.
        </p>
      </div>
      <div className="p-8">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Item Name <span className="text-red-500">*</span>
            </label>
            <Input
              required
              value={formData.name}
              onChange={(e) =>
                setFormData((p) => ({ ...p, name: e.target.value }))
              }
              placeholder="e.g. Aloo Tikki Burger"
              className="focus-visible:ring-red-500 w-full"
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Description
            </label>
            <Textarea
              value={formData.description}
              onChange={(e) =>
                setFormData((p) => ({ ...p, description: e.target.value }))
              }
              placeholder="Brief description of the dish..."
              className="focus-visible:ring-red-500 resize-none w-full"
              rows={3}
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Price (₹) <span className="text-red-500">*</span>
            </label>
            <Input
              type="number"
              step="0.01"
              required
              value={formData.price}
              onChange={(e) =>
                setFormData((p) => ({ ...p, price: e.target.value }))
              }
              placeholder="99.00"
              className="focus-visible:ring-red-500 w-full"
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Item Image
            </label>
            <Input
              type="file"
              accept="image/*"
              onChange={(e) => setImageFile(e.target.files?.[0] || null)}
              className="cursor-pointer file:text-gray-600 focus-visible:ring-red-500 w-full"
            />
          </div>

          <div className="flex items-center justify-between p-4 border border-gray-200 rounded-xl bg-gray-50/50 mt-4">
            <div>
              <label className="text-base font-bold text-gray-900">
                Available to Order
              </label>
              <p className="text-sm text-gray-500">
                Turn this off if the item is out of stock.
              </p>
            </div>
            <Switch
              checked={formData.is_available}
              onCheckedChange={(checked) =>
                setFormData((p) => ({ ...p, is_available: checked }))
              }
              className="data-[state=checked]:bg-green-500"
            />
          </div>

          <div className="pt-4 border-t border-gray-100">
            <Button
              type="submit"
              disabled={loading}
              className="w-full bg-red-600 hover:bg-red-700 text-white py-6 text-lg font-semibold rounded-xl transition-all shadow-md"
            >
              {loading ? "Adding Item..." : "Save Menu Item"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
