import React from "react";

export function CreateRestaurantPrompt() {
  return (
    <div className="max-w-2xl mx-auto mt-12 bg-white border border-gray-200 rounded-xl shadow-sm p-10 text-center">
      <div className="w-16 h-16 bg-red-50 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">
        🍽️
      </div>
      <h2 className="text-2xl font-bold text-gray-800 mb-2">
        No Restaurant Found
      </h2>
      <p className="text-gray-500 mb-8">
        You haven't set up a restaurant profile yet. Get started to manage your
        business.
      </p>
      <button
        onClick={() => {
          /* Add your navigation or modal trigger here */
          console.log("Create Restaurant clicked");
        }}
        className="px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors shadow-md hover:shadow-lg"
      >
        + Create Restaurant
      </button>
    </div>
  );
}
