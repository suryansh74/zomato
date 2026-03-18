import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./Pages/Home";
import Login from "./Pages/Login";
import SelectRole from "./Pages/SelectRole";
import AuthProvider from "./context/AuthProvider";
import ProtectedRoute from "./middleware/ProtectedRoute";
import RoleProtectedRoute from "./middleware/RoleProtectedRoute";
import PublicRoute from "./middleware/PublicRoute";
import { Toaster } from "react-hot-toast";
import { ThemeProvider } from "./components/theme-provider";

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            {/* public — anyone */}
            <Route element={<PublicRoute />}>
              <Route path="/login" element={<Login />} />
            </Route>

            {/* logged in but no role yet */}
            <Route element={<ProtectedRoute />}>
              <Route path="/select-role" element={<SelectRole />} />
            </Route>

            {/* logged in AND has role */}
            <Route element={<RoleProtectedRoute />}>
              <Route path="/" element={<Home />} />
            </Route>
          </Routes>

          <Toaster />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}
