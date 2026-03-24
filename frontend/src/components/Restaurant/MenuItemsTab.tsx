import React, { useState } from "react";
import axios from "axios";
import toast from "react-hot-toast";
import { restaurantServiceUrl } from "@/lib/config";
import type { MenuItem } from "@/types/types";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";

interface MenuItemsTabProps {
  menuItems: MenuItem[];
  loading: boolean;
  onAddClick: () => void;
  onRefresh: () => void;
}

export function MenuItemsTab({
  menuItems,
  loading,
  onAddClick,
  onRefresh,
}: MenuItemsTabProps) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden animate-in fade-in duration-300">
      <div className="p-6 border-b border-gray-100 bg-gray-50/50 flex justify-between items-center">
        <div>
          <h3 className="text-xl font-bold text-gray-900">Your Menu</h3>
          <p className="text-sm text-gray-500">
            Manage your existing menu items and availability.
          </p>
        </div>
        {/* Quick Add Button in Header */}
        {menuItems.length > 0 && (
          <Button
            onClick={onAddClick}
            className="bg-red-50 text-red-600 hover:bg-red-100 font-semibold shadow-none border-0"
          >
            + Add New Item
          </Button>
        )}
      </div>

      <div className="p-6">
        {loading ? (
          <div className="flex justify-center py-12">
            <span className="animate-pulse text-gray-400 font-medium">
              Loading your delicious menu...
            </span>
          </div>
        ) : menuItems.length === 0 ? (
          <div className="text-center py-16 px-4">
            <div className="text-5xl mb-4">🍽️</div>
            <h4 className="text-lg font-bold text-gray-900 mb-2">
              Your menu is empty
            </h4>
            <p className="text-gray-500 mb-6 max-w-sm mx-auto">
              You haven't added any dishes yet. Let's get started by adding your
              first menu item.
            </p>
            <Button
              onClick={onAddClick}
              className="bg-red-600 hover:bg-red-700 text-white px-8 py-2 rounded-lg"
            >
              + Add Your First Item
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {menuItems.map((item) => (
              <MenuItemCard key={item.id} item={item} onRefresh={onRefresh} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

// ============================================================================
// SUB-COMPONENT: SINGLE MENU ITEM CARD (VIEW ONLY + EDIT BUTTON)
// ============================================================================
function MenuItemCard({
  item,
  onRefresh,
}: {
  item: MenuItem;
  onRefresh: () => void;
}) {
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Keep quick toggle on the card for fast UX
  const toggleAvailability = async (checked: boolean) => {
    try {
      const payload = new FormData();
      payload.append("is_available", String(checked));
      await axios.put(`${restaurantServiceUrl}/menu/${item.id}`, payload, {
        withCredentials: true,
      });
      toast.success(
        checked ? "Item marked Available" : "Item marked Out of Stock",
      );
      onRefresh();
    } catch (err) {
      toast.error("Failed to update status");
    }
  };

  return (
    <>
      <div className="flex border border-gray-200 rounded-2xl overflow-hidden shadow-sm bg-white hover:shadow-md transition-shadow h-44 group w-full relative">
        {/* Image Side */}
        <div className="w-36 bg-gray-100 flex-shrink-0 relative border-r border-gray-100">
          {item.image ? (
            <img
              src={item.image}
              alt={item.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex flex-col items-center justify-center text-gray-400 p-2 bg-gray-50">
              <span className="text-2xl mb-1">🍲</span>
              <span className="text-[10px] uppercase font-bold tracking-wider">
                No Image
              </span>
            </div>
          )}
          {!item.is_available && (
            <div className="absolute inset-0 bg-white/70 flex items-center justify-center backdrop-blur-sm">
              <span className="bg-red-600 text-white text-[10px] uppercase font-extrabold tracking-wider px-3 py-1.5 rounded-full shadow-sm">
                Out of Stock
              </span>
            </div>
          )}
        </div>

        {/* Content Side */}
        <div className="p-4 flex flex-col flex-1 justify-between relative overflow-hidden">
          <div>
            <div className="flex justify-between items-start gap-2">
              <h4 className="font-bold text-gray-900 text-lg leading-tight line-clamp-1 truncate">
                {item.name}
              </h4>
              <span className="font-bold text-green-700 bg-green-50 px-2 py-0.5 rounded text-sm whitespace-nowrap">
                ₹{item.price}
              </span>
            </div>
            <p className="text-sm text-gray-500 line-clamp-2 mt-1.5 leading-snug">
              {item.description}
            </p>
          </div>

          <div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-100">
            <div
              className="flex items-center gap-2 cursor-pointer"
              onClick={() => toggleAvailability(!item.is_available)}
            >
              <Switch
                checked={item.is_available}
                className="data-[state=checked]:bg-green-500 pointer-events-none"
              />
              <span className="text-xs font-semibold text-gray-600 select-none">
                {item.is_available ? "In Stock" : "Out of Stock"}
              </span>
            </div>

            {/* STYLISH EDIT BUTTON */}
            <Button
              onClick={() => setIsModalOpen(true)}
              size="sm"
              className="bg-gray-100 hover:bg-blue-600 text-gray-700 hover:text-white transition-colors duration-300 font-semibold px-4 flex items-center gap-1.5 shadow-none"
            >
              <svg
                className="w-4 h-4"
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
              Edit
            </Button>
          </div>
        </div>
      </div>

      {/* RENDER MODAL IF OPEN */}
      {isModalOpen && (
        <EditMenuItemModal
          item={item}
          onClose={() => setIsModalOpen(false)}
          onRefresh={onRefresh}
        />
      )}
    </>
  );
}

// ============================================================================
// SUB-COMPONENT: EDIT MODAL POP-UP
// ============================================================================
function EditMenuItemModal({
  item,
  onClose,
  onRefresh,
}: {
  item: MenuItem;
  onClose: () => void;
  onRefresh: () => void;
}) {
  const [loading, setLoading] = useState(false);
  const [editData, setEditData] = useState({
    name: item.name,
    description: item.description || "",
    price: String(item.price),
    is_available: item.is_available,
  });
  const [newImage, setNewImage] = useState<File | null>(null);

  const handleUpdate = async () => {
    try {
      setLoading(true);
      const payload = new FormData();
      payload.append("name", editData.name);
      payload.append("description", editData.description);
      payload.append("price", editData.price);
      payload.append("is_available", String(editData.is_available));
      if (newImage) payload.append("image", newImage);

      await axios.put(`${restaurantServiceUrl}/menu/${item.id}`, payload, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      });

      toast.success("Item updated successfully");
      onRefresh();
      onClose();
    } catch (err) {
      toast.error("Failed to update item");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (
      !confirm(
        `Are you sure you want to completely delete "${item.name}"? This cannot be undone.`,
      )
    )
      return;
    try {
      setLoading(true);
      await axios.delete(`${restaurantServiceUrl}/menu/${item.id}`, {
        withCredentials: true,
      });
      toast.success("Item deleted permanently");
      onRefresh();
      onClose(); // Close modal on delete
    } catch (err) {
      toast.error("Failed to delete item");
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-gray-900/60 backdrop-blur-sm animate-in fade-in duration-200">
      {/* Modal Container */}
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
          <h3 className="text-xl font-bold text-gray-900">Edit Menu Item</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors p-1"
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
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Scrollable Form Body */}
        <div className="p-6 overflow-y-auto space-y-5">
          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Item Name
            </label>
            <Input
              value={editData.name}
              onChange={(e) =>
                setEditData((p) => ({ ...p, name: e.target.value }))
              }
              className="focus-visible:ring-blue-500 w-full"
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Description
            </label>
            <Textarea
              value={editData.description}
              onChange={(e) =>
                setEditData((p) => ({ ...p, description: e.target.value }))
              }
              rows={3}
              className="focus-visible:ring-blue-500 resize-none w-full"
            />
          </div>

          <div className="flex gap-4">
            <div className="w-1/2">
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">
                Price (₹)
              </label>
              <Input
                type="number"
                step="0.01"
                value={editData.price}
                onChange={(e) =>
                  setEditData((p) => ({ ...p, price: e.target.value }))
                }
                className="focus-visible:ring-blue-500 w-full"
              />
            </div>
            <div className="w-1/2">
              <label className="block text-sm font-semibold text-gray-700 mb-1.5">
                Availability
              </label>
              <div className="flex items-center justify-between border border-gray-200 rounded-md px-4 h-10 bg-gray-50">
                <span className="text-sm font-medium text-gray-700">
                  {editData.is_available ? "In Stock" : "Out of Stock"}
                </span>
                <Switch
                  checked={editData.is_available}
                  onCheckedChange={(c) =>
                    setEditData((p) => ({ ...p, is_available: c }))
                  }
                  className="data-[state=checked]:bg-green-500"
                />
              </div>
            </div>
          </div>

          <div>
            <label className="block text-sm font-semibold text-gray-700 mb-1.5">
              Update Image (Optional)
            </label>
            <Input
              type="file"
              accept="image/*"
              onChange={(e) => setNewImage(e.target.files?.[0] || null)}
              className="cursor-pointer file:text-gray-600 focus-visible:ring-blue-500 w-full text-sm"
            />
          </div>
        </div>

        {/* Footer Actions (Delete on Left, Cancel/Save on Right) */}
        <div className="p-4 border-t border-gray-100 bg-gray-50 flex items-center justify-between">
          <Button
            onClick={handleDelete}
            disabled={loading}
            variant="ghost"
            className="text-red-600 hover:bg-red-50 hover:text-red-700 px-3"
          >
            Delete Item
          </Button>

          <div className="flex gap-3">
            <Button
              onClick={onClose}
              variant="outline"
              disabled={loading}
              className="bg-white"
            >
              Cancel
            </Button>
            <Button
              onClick={handleUpdate}
              disabled={loading}
              className="bg-blue-600 hover:bg-blue-700 text-white min-w-[120px]"
            >
              {loading ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
