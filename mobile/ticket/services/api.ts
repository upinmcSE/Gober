import { Platform } from "react-native";
import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";
import AsyncStorage from "@react-native-async-storage/async-storage";

const url = "http://192.168.48.230:8082";

const x_api_key = "894f5978-f86b-4b04-9bc0-328e7ac6166e"

const Api: AxiosInstance = axios.create({ baseURL: url + "/v1/2025" });

Api.interceptors.request.use(async config => {
    const token = await AsyncStorage.getItem("token");

    if (token) config.headers.set("Authorization", `Bearer ${token}`);

    if (x_api_key) config.headers.set("X-API-Key", x_api_key);

    return config;
});

Api.interceptors.response.use(
    async (res: AxiosResponse) => res.data,
    async (err: AxiosError) => Promise.reject(err)
);

export { Api };