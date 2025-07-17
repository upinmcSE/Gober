import AsyncStorage from "@react-native-async-storage/async-storage";
import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";

const url = "http://172.20.10.4:8082";

const x_api_key = "894f5978-f86b-4b04-9bc0-328e7ac6166e"

const Api: AxiosInstance = axios.create({ 
    baseURL: url + "/v1/2025", 
    timeout: 10000,
});

Api.interceptors.request.use(async config => {
    const token = await AsyncStorage.getItem("token");
    console.log("token", token)

    config.headers.set("Authorization", `Bearer ${token}`);

    if (x_api_key) config.headers.set("X-API-Key", x_api_key);

    console.log("Final headers:", config.headers);
    return config;
});

Api.interceptors.response.use(
    async (res: AxiosResponse) => res.data,
    async (err: AxiosError) => Promise.reject(err)
);

export { Api };

