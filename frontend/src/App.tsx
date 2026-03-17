import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./Pages/Home";
import Login from "./Pages/Login";
import Dashboard from "./Pages/Dashboard";
import SelectRole from "./Pages/SelectRole";
import AuthProvider from "./context/AuthProvider";
import ProtectedRoute from "./middleware/ProtectedRoute";
import RoleProtectedRoute from "./middleware/RoleProtectedRoute";
import PublicRoute from "./middleware/PublicRoute";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* public — anyone */}
          <Route path="/" element={<Home />} />
          <Route element={<PublicRoute />}>
            <Route path="/login" element={<Login />} />
          </Route>

          {/* logged in but no role yet */}
          <Route element={<ProtectedRoute />}>
            <Route path="/select-role" element={<SelectRole />} />
          </Route>

          {/* logged in AND has role */}
          <Route element={<RoleProtectedRoute />}>
            <Route path="/dashboard" element={<Dashboard />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
