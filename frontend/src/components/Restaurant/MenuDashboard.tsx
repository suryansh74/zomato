import React, { useState, useEffect } from "react";
import axios from "axios";
import toast from "react-hot-toast";
import { restaurantServiceUrl } from "@/lib/config";
import type { MenuItem } from "@/types/types";

// Import our new sub-components
import { MenuItemsTab } from "./MenuItemsTab";
import { AddItemTab } from "./AddItemTab";
import { SalesTab } from "./SalesTab";

export function MenuDashboard() {
  const [activeTab, setActiveTab] = useState<
    "menu-items" | "add-item" | "sales"
  >("menu-items");
  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchMenuItems = async () => {
    try {
      setLoading(true);
      const res = await axios.get(`${restaurantServiceUrl}/menu`, {
        withCredentials: true,
      });

      // FIXED: Safely extracting the nested menu_items array based on your Go backend response
      const responseData = res.data?.data;
      if (responseData?.menu_items && Array.isArray(responseData.menu_items)) {
        setMenuItems(responseData.menu_items);
      } else if (Array.isArray(responseData)) {
        // Fallback just in case you update your Go backend later to remove the "menu_items" wrapper
        setMenuItems(responseData);
      } else {
        setMenuItems([]);
      }
    } catch (err) {
      toast.error("Failed to load menu items");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMenuItems();
  }, []);

  return (
    <div className="w-full mt-12 pb-20">
      {/* CUSTOM TAB HEADER */}
      <div className="flex w-full border-b border-gray-200 mb-8">
        <button
          onClick={() => setActiveTab("menu-items")}
          className={`flex-1 py-4 text-center text-base font-semibold transition-all duration-200 ${
            activeTab === "menu-items"
              ? "text-red-600 border-b-2 border-red-600"
              : "text-gray-500 border-b-2 border-transparent hover:text-gray-700 hover:border-gray-300"
          }`}
        >
          Menu Items
        </button>
        <button
          onClick={() => setActiveTab("add-item")}
          className={`flex-1 py-4 text-center text-base font-semibold transition-all duration-200 ${
            activeTab === "add-item"
              ? "text-red-600 border-b-2 border-red-600"
              : "text-gray-500 border-b-2 border-transparent hover:text-gray-700 hover:border-gray-300"
          }`}
        >
          Add New Item
        </button>
        <button
          onClick={() => setActiveTab("sales")}
          className={`flex-1 py-4 text-center text-base font-semibold transition-all duration-200 ${
            activeTab === "sales"
              ? "text-red-600 border-b-2 border-red-600"
              : "text-gray-500 border-b-2 border-transparent hover:text-gray-700 hover:border-gray-300"
          }`}
        >
          Sales Overview
        </button>
      </div>

      {/* TAB CONTENT AREA */}
      <div className="w-full">
        {activeTab === "menu-items" && (
          <MenuItemsTab
            menuItems={menuItems}
            loading={loading}
            onAddClick={() => setActiveTab("add-item")}
            onRefresh={fetchMenuItems}
          />
        )}

        {activeTab === "add-item" && (
          <AddItemTab
            onSuccess={() => {
              fetchMenuItems();
              setActiveTab("menu-items");
            }}
          />
        )}

        {activeTab === "sales" && <SalesTab />}
      </div>
    </div>
  );
}
