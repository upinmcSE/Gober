import AsyncStorage from "@react-native-async-storage/async-storage";
import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";

const url = "http://192.168.1.6:8082";

const x_api_key = "894f5978-f86b-4b04-9bc0-328e7ac6166e"

const Api: AxiosInstance = axios.create({ 
    baseURL: url + "/v1/2025", 
    timeout: 10000,
});

Api.interceptors.request.use(async config => {
    const token = await AsyncStorage.getItem("accessToken");
    console.log("accessToken", token)

    if (token) config.headers['Authorization'] = `Bearer ${token}`;

    if (x_api_key) config.headers.set("X-API-Key", x_api_key);

    //console.log("Final headers:", config.headers);
    return config;
});

Api.interceptors.response.use(
    async (res: AxiosResponse) => res.data,
    async (err: AxiosError) => {
        console.log("Response Error:", err.response?.status);
        console.log("Response Headers:", err.response?.headers);
        console.log("Response Data:", err.response?.data);
        console.log("Request Headers:", err.config?.headers);
        return Promise.reject(err)}
);

export { Api };

