import { create } from 'zustand';

interface User {
  id?: string;
  email: string;
  name: string;
  token?: string;
}

interface AuthState {
  user: User | null;
  isInitialized: boolean;
  setUser: (user: User | null) => void;
  login: (user: User) => void;
  initialize: () => Promise<void>;
  logout: () => void;
}

export const useAuth = create<AuthState>((set) => ({
  user: null,
  isInitialized: false,
  setUser: (user) => set({ user }),
  login: (user) => set({ user }), 
  initialize: async () => {
    try {
      const response = await fetch('http://localhost:8000/api/user', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include', 
      });

      if (!response.ok) {
        throw new Error('Unauthorized');
      }

      const user = await response.json();
      set({ user, isInitialized: true });
    } catch (error) {
      console.error('Failed to fetch user:', error);
      set({ user: null, isInitialized: true });
    }
  },
  logout: () => {
    set({ user: null });
    localStorage.removeItem('jwt'); 
  },
}));
