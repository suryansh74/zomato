import { createContext } from "react";
import type { CartContextType } from "../types/types";

export const CartContext = createContext<CartContextType>({
  cartLength: 0,
  cartItems: [],
  subtotal: 0,
  isCartLoading: true, // ✅ Default to true so it spins immediately on reload
  setCartLength: () => {},
  fetchCart: async () => {},
});
