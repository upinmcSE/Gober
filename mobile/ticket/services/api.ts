import { Platform } from "react-native";
import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";
import AsyncStorage from "@react-native-async-storage/async-storage";

const url = Platform.OS === "android" ? "http://192.168.1.12:8080" : "http://localhost:8080";

const x_api_key = "Ft89HodRC7Sfvy85MfH9xhXPRHiBVqVNCAeQ30rsqv0nxvsxcFhOCZRxgMamERMn"

const Api: AxiosInstance = axios.create({ baseURL: url + "/v1/2025" });

Api.interceptors.request.use(async config => {
    const token = await AsyncStorage.getItem("token");

    if (token) config.headers.set("Authorization", `Bearer ${token}`);

    if (x_api_key) config.headers.set("x-api-key", x_api_key);

    return config;
});

Api.interceptors.response.use(
    async (res: AxiosResponse) => res.data,
    async (err: AxiosError) => Promise.reject(err)
);

export { Api };