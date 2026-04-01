import { BrowserRouter, Routes, Route, useLocation } from "react-router-dom";
import Home from "./Pages/Home";
import Login from "./Pages/Login";
import SelectRole from "./Pages/SelectRole";
import AuthProvider from "./context/AuthProvider";
import ProtectedRoute from "./middleware/ProtectedRoute";
import RoleProtectedRoute from "./middleware/RoleProtectedRoute";
import PublicRoute from "./middleware/PublicRoute";
import { Toaster } from "react-hot-toast";
import Account from "./Pages/Account";

import Navbar from "./components/Navbar";
import { useAuth } from "./context/useAuth";
import Restaurant from "./Pages/Restaurant";
import CustomerRestaurant from "./components/Restaurant/CustomerRestaurant";
import CartProvider from "./context/CartProvider";
import Cart from "./components/Restaurant/Cart";
import AddAddressPage from "./Pages/Address";
import Success from "./Pages/Success";

// separate component so useLocation works inside BrowserRouter
function Layout() {
  const location = useLocation();
  const { user } = useAuth();

  const hideNavbar =
    ["/login", "/select-role"].includes(location.pathname) ||
    user?.role === "restaurant_owner";

  return (
    <>
      {!hideNavbar && <Navbar />}

      <Routes>
        <Route element={<PublicRoute />}>
          <Route path="/login" element={<Login />} />
        </Route>

        <Route element={<ProtectedRoute />}>
          <Route path="/select-role" element={<SelectRole />} />
          <Route path="/account" element={<Account />} />
        </Route>

        <Route element={<RoleProtectedRoute />}>
          <Route
            path="/"
            element={
              user?.role === "restaurant_owner" ? <Restaurant /> : <Home />
            }
          />
          <Route path="/restaurant/:id" element={<CustomerRestaurant />} />
          <Route path="/address" element={<AddAddressPage />} />
          <Route path="/cart" element={<Cart />} />
          <Route path="/success" element={<Success />} />
        </Route>
      </Routes>

      <Toaster />
    </>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <CartProvider>
        <BrowserRouter>
          <Layout />
        </BrowserRouter>
      </CartProvider>
    </AuthProvider>
  );
}
