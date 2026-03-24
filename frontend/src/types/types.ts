export interface User {
  name: string;
  email: string;
  image: string;
  role: string;
}

export interface AuthContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  loading: boolean;
  location: LocationData | null;
  loadingLocation: boolean;
  city: string;
}

export interface LocationData {
  latitude: number;
  longitude: number;
  formattedAddress: string;
}

export type Role = "customer" | "restaurant_owner" | "rider" | null;

export interface Restaurant {
  id: string;
  name: string;
  description?: string;
  image: string;
  owner_email: string;
  phone: string;
  is_verified: boolean;
  auto_location: {
    type: "Point";
    coordinates: [number, number]; // [longitude, latitude]
    formatted_address: string;
  };
  is_open: boolean;
  created_at: string;
  updated_at: string;
}
