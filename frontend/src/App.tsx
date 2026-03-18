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

// separate component so useLocation works inside BrowserRouter
function Layout() {
  const location = useLocation();
  const hideNavbar = ["/login", "/select-role"].includes(location.pathname);

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
          <Route path="/" element={<Home />} />
        </Route>
      </Routes>
      <Toaster />
    </>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Layout />
      </BrowserRouter>
    </AuthProvider>
  );
}
