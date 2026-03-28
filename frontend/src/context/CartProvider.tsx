import { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { CartContext } from "./CartContext";
import { useAuth } from "./useAuth";
import { restaurantServiceUrl } from "@/lib/config";
import type { PopulatedCartItem } from "@/types/types";

export default function CartProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [cartLength, setCartLength] = useState(0);
  const [cartItems, setCartItems] = useState<PopulatedCartItem[]>([]);
  const [subtotal, setSubtotal] = useState(0);
  const [isCartLoading, setIsCartLoading] = useState(true);

  // ✅ 1. Grab 'loading' from useAuth and rename it to 'authLoading' for clarity
  const { user, loading: authLoading } = useAuth();

  const fetchCart = useCallback(async () => {
    // ✅ 2. If Auth is still figuring out who the user is, DO NOT turn off the spinner!
    if (authLoading) {
      return;
    }

    // ✅ 3. If Auth is finished, but there is no user, clear the cart.
    if (!user) {
      setCartLength(0);
      setCartItems([]);
      setSubtotal(0);
      setIsCartLoading(false);
      return;
    }

    try {
      const res = await axios.get(`${restaurantServiceUrl}/cart`, {
        withCredentials: true,
      });

      const items = res.data?.data?.cart || res.data?.cart || [];
      setCartItems(items);
      setCartLength(res.data?.data?.cartLength || res.data?.cartLength || 0);
      setSubtotal(res.data?.data?.subtotal || res.data?.subtotal || 0);
    } catch (error) {
      console.error("Failed to fetch cart:", error);
      setCartLength(0);
      setCartItems([]);
      setSubtotal(0);
    } finally {
      setIsCartLoading(false); // Turn off spinner after the fetch completes
    }
  }, [user, authLoading]); // ✅ 4. Add authLoading to the dependency array

  useEffect(() => {
    fetchCart();
  }, [fetchCart]);

  return (
    <CartContext.Provider
      value={{
        cartLength,
        cartItems,
        subtotal,
        isCartLoading,
        setCartLength,
        fetchCart,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}
