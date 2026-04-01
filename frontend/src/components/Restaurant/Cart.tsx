import { useCart } from "@/context/useCart";
import { restaurantServiceUrl } from "@/lib/config";
import axios from "axios";
import { Minus, Plus, ShoppingBag, Trash2, Loader2 } from "lucide-react";
import toast from "react-hot-toast";
import { Link, useNavigate } from "react-router-dom"; // ✅ Imported useNavigate

export default function Cart() {
  const { cartItems, subtotal, fetchCart, isCartLoading } = useCart();
  const navigate = useNavigate(); // ✅ Initialized navigate

  const handleUpdateQuantity = async (
    itemId: string,
    action: "inc" | "dec",
  ) => {
    try {
      await axios.put(
        `${restaurantServiceUrl}/cart/update`,
        { itemId, action },
        { withCredentials: true },
      );
      // Refresh global state to show new numbers
      fetchCart();
    } catch (error) {
      toast.error("Failed to update quantity");
    }
  };

  const handleClearCart = async () => {
    try {
      await axios.delete(`${restaurantServiceUrl}/cart/clear`, {
        withCredentials: true,
      });
      toast.success("Cart cleared");
      fetchCart();
    } catch (error) {
      toast.error("Failed to clear cart");
    }
  };

  if (isCartLoading) {
    return (
      <div className="min-h-[70vh] flex flex-col items-center justify-center bg-gray-50">
        <Loader2 className="w-12 h-12 text-red-500 animate-spin mb-4" />
        <h2 className="text-xl font-semibold text-gray-700">
          Loading your cart...
        </h2>
      </div>
    );
  }

  if (cartItems.length === 0) {
    return (
      <div className="min-h-[70vh] flex flex-col items-center justify-center bg-gray-50">
        <ShoppingBag className="w-24 h-24 text-gray-300 mb-6" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Your cart is empty
        </h2>
        <p className="text-gray-500 mb-8">
          You can go to home page to view more restaurants
        </p>
        <Link
          to="/"
          className="bg-red-500 text-white px-6 py-3 rounded-lg font-bold hover:bg-red-600 transition"
        >
          See Restaurants near you
        </Link>
      </div>
    );
  }

  // Get the restaurant info from the first item
  const restaurant = cartItems[0].restaurantId;

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-5xl mx-auto flex flex-col lg:flex-row gap-8">
        {/* LEFT COLUMN: Cart Items */}
        <div className="flex-1 space-y-6">
          <div className="flex items-center justify-between bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
            <div>
              <h2 className="text-xl font-bold text-gray-900">
                {restaurant.name}
              </h2>
              <p className="text-sm text-gray-500">
                {restaurant.auto_location.formatted_address}
              </p>
            </div>
            <button
              onClick={handleClearCart}
              className="text-gray-400 hover:text-red-500 transition flex items-center gap-2 text-sm font-medium"
            >
              <Trash2 className="w-4 h-4" /> Clear Cart
            </button>
          </div>

          <div className="space-y-4">
            {cartItems.map((cartItem) => (
              <div
                key={cartItem._id}
                className="flex items-center bg-white p-4 rounded-xl shadow-sm border border-gray-100"
              >
                {/* Item Image */}
                <div className="w-20 h-20 bg-gray-100 rounded-lg overflow-hidden shrink-0">
                  {cartItem.itemId.image ? (
                    <img
                      src={cartItem.itemId.image}
                      alt={cartItem.itemId.name}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-xs text-gray-400">
                      No Image
                    </div>
                  )}
                </div>

                {/* Item Details */}
                <div className="ml-4 flex-1">
                  <h3 className="font-bold text-gray-900">
                    {cartItem.itemId.name}
                  </h3>
                  <p className="text-gray-500 mt-1">₹{cartItem.itemId.price}</p>
                </div>

                {/* Quantity Controls */}
                <div className="flex items-center gap-4 ml-4">
                  <button
                    onClick={() =>
                      handleUpdateQuantity(cartItem.itemId.id, "dec")
                    }
                    className="w-8 h-8 rounded-full border border-gray-300 flex items-center justify-center text-gray-600 hover:bg-gray-50 hover:border-gray-400 transition"
                  >
                    <Minus className="w-4 h-4" />
                  </button>
                  <span className="font-semibold text-gray-900 w-4 text-center">
                    {cartItem.quantity}
                  </span>
                  <button
                    onClick={() =>
                      handleUpdateQuantity(cartItem.itemId.id, "inc")
                    }
                    className="w-8 h-8 rounded-full border border-red-200 flex items-center justify-center text-red-500 hover:bg-red-50 transition"
                  >
                    <Plus className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* RIGHT COLUMN: Bill Details */}
        <div className="w-full lg:w-[350px]">
          <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 sticky top-24">
            <h3 className="font-bold text-lg text-gray-900 mb-6">
              Bill Details
            </h3>

            <div className="space-y-4 text-sm text-gray-600 mb-6">
              <div className="flex justify-between">
                <span>Item Total</span>
                <span className="font-medium">₹{subtotal}</span>
              </div>
              <div className="flex justify-between">
                <span>Delivery Fee</span>
                <span className="font-medium text-green-600">FREE</span>
              </div>
              <div className="flex justify-between border-b border-dashed pb-4">
                <span>Platform Fee</span>
                <span className="font-medium">₹5</span>
              </div>
              <div className="flex justify-between text-lg font-bold text-gray-900 pt-2">
                <span>To Pay</span>
                <span>₹{subtotal + 5}</span>
              </div>
            </div>

            <button
              disabled={!restaurant.is_open}
              onClick={() => navigate("/address")} // ✅ Added the onClick routing
              className={`w-full font-bold py-4 rounded-xl transition shadow-md mt-6 ${
                restaurant.is_open
                  ? "bg-red-500 text-white hover:bg-red-600 shadow-red-200"
                  : "bg-gray-300 text-gray-500 cursor-not-allowed shadow-none"
              }`}
            >
              {restaurant.is_open ? "Proceed to Checkout" : "Restaurant Closed"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
