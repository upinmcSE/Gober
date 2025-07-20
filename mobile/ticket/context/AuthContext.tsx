import { userService } from '@/services/user';
import { User } from '@/types/user';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { router } from 'expo-router';
import React, { useEffect, useState } from 'react';

interface AuthContextProps {
  isLoggedIn: boolean;
  isLoadingAuth: boolean;
  authenticate: (authMode: "login" | "register", email: string, password: string) => Promise<void>;
  logout: VoidFunction;
  user: User | null;
}

const AuthContext = React.createContext({} as AuthContextProps);

export function useAuth() {
  return React.useContext(AuthContext);
}

export function AuthenticationProvider({ children }: React.PropsWithChildren) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isLoadingAuth, setIsLoadingAuth] = useState(false);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    async function checkIfLoggedIn() {
      const token = await AsyncStorage.getItem('accessToken');
      const user = await AsyncStorage.getItem('user');

      //console.log("Checking if user is logged in", { token, user });

      if (token && user) {
        setIsLoggedIn(true);
        setUser(JSON.parse(user));
        // Redirect to authed route
        router.replace('/(authed)/(tabs)/(events)');
      }
    }

    checkIfLoggedIn();
  }, []);

  async function authenticate(authMode: "login" | "register", email: string, password: string) {
  console.log("but but");
  try {
    setIsLoadingAuth(true);

    const response = await userService[authMode](email, password);
    console.log("Authentication response:", response);

    
    const data = response?.data[0];

    if (data) {
      const accessToken = data.access_token;
      const refreshToken = data.refresh_token
      const user = data.of_account;

      if (accessToken && user) {
        setIsLoggedIn(true);
        setUser(user);

        await AsyncStorage.setItem('accessToken', accessToken);
        await AsyncStorage.setItem('refreshToken', refreshToken);
        await AsyncStorage.setItem('user', JSON.stringify(user));

        router.replace("/(authed)/(tabs)/(events)");
      } else {
        console.error("Missing token or user in response", data);
      }
    } else {
      console.error("No data returned from server");
    }
  } catch (error) {
    console.error("Authentication error:", error);
    setIsLoggedIn(false);
  } finally {
    setIsLoadingAuth(false);
  }
}

  async function logout() {
    setIsLoggedIn(false);
    setUser(null);

    await AsyncStorage.removeItem('accessToken');
    await AsyncStorage.removeItem('user');
  }

  return (
    <AuthContext.Provider
      value={ {
        authenticate,
        logout,
        isLoggedIn,
        isLoadingAuth,
        user
      } }>
      { children }
    </AuthContext.Provider>
  );
}