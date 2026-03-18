import { useAuth } from "@/context/useAuth";
import { MapPin, Search, ShoppingCart } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useLocation, useSearchParams } from "react-router-dom";

export default function Navbar() {
  const location = useLocation();
  const isHome = location.pathname === "/";

  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState(searchParams.get("search") || "");

  const { city } = useAuth();

  useEffect(() => {
    if (!isHome) return;

    const timer = setTimeout(() => {
      if (search) {
        setSearchParams({ search });
      } else {
        setSearchParams({});
      }
    }, 400);

    return () => clearTimeout(timer);
  }, [search, isHome, setSearchParams]);

  return (
    <div className="w-full bg-white shadow-sm sticky top-0 z-50">
      {/* Top Navbar */}
      <div className="border-b">
        <div className="max-w-7xl mx-auto flex items-center justify-between px-4 py-3">
          <Link
            to="/"
            className="text-red-500 text-2xl font-bold tracking-wide"
          >
            Zomato
          </Link>

          <div className="flex items-center gap-6">
            <Link to="/cart" className="relative group">
              <ShoppingCart className="text-gray-700 w-6 h-6 group-hover:text-red-500 transition" />
              <span className="absolute -top-2 -right-2 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-xs font-bold text-white">
                0
              </span>
            </Link>
            <Link
              to="/account"
              className="text-gray-700 text-sm hover:text-red-500 transition"
            >
              Account
            </Link>
          </div>
        </div>
      </div>

      {/* Search Section */}
      {isHome && (
        <div className="bg-gray-50 py-4">
          <div className="max-w-4xl mx-auto px-4">
            <div className="flex items-center bg-white rounded-xl shadow-sm overflow-hidden">
              {/* Location */}
              <div className="flex items-center gap-2 px-4 border-r min-w-[180px]">
                <MapPin className="text-red-500 w-5 h-5" />
                <span className="text-gray-600 text-sm truncate">{city}</span>
              </div>

              {/* Search */}
              <div className="flex items-center flex-1 px-4">
                <Search className="text-gray-400 w-5 h-5 mr-2" />
                <input
                  type="text"
                  placeholder="Search for restaurant, cuisine or dish"
                  className="w-full py-3 text-sm outline-none placeholder-gray-400"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
