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
