import React from "react";

export function SalesTab() {
  return (
    <div className="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden animate-in fade-in duration-300">
      <div className="p-6 border-b border-gray-100 bg-gray-50/50">
        <h3 className="text-xl font-bold text-gray-900">Sales Dashboard</h3>
        <p className="text-sm text-gray-500">
          View your restaurant's performance and analytics.
        </p>
      </div>
      <div className="p-12">
        <div className="h-64 flex flex-col items-center justify-center border-2 border-dashed border-gray-200 rounded-xl bg-gray-50 text-gray-400">
          <span className="text-4xl mb-3">📈</span>
          <p className="font-medium">Sales analytics coming soon...</p>
        </div>
      </div>
    </div>
  );
}
