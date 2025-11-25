export interface RegisterRequest {
  email: string;
  password: string;
}

export interface RegisterResponse {
  message: string;
}

export interface ErrorResponse {
  error: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
  public_id: string;
  email: string;
  role: string;
  created_at: string;
  display_name?: string;
}

export interface LoginResponse {
  message: string;
  access_token: string;
  user: User;
}

export interface RefreshResponse {
  message: string;
  access_token: string;
  user: User;
}

// Additional types for better type safety
export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface AuthContextType {
  user: User | null;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
  isAuthenticated: boolean;
  isLoading: boolean;
}

export interface ApiError {
  message: string;
  status?: number;
}
