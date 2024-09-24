import {fetchBaseQuery, retry} from "@reduxjs/toolkit/query/react";
import {BASE_URL} from "../../constants";
import type {RootState} from "../store";
import {createApi} from "@reduxjs/toolkit/query/react";

const baseQuery = fetchBaseQuery({
    baseUrl: `${BASE_URL}/`,
    // prepareHeaders: (headers, {getState}) => {
    //     const token = (getState() as RootState).auth.token || localStorage.getItem("token");
    //     if (token) {
    //         headers.set("Authorization", `Bearer ${token}`);
    //     }
    //
    //     return headers;
    // }
})

const baseQueryWithRetry = retry(baseQuery, {maxRetries: 1})

export const api = createApi({
    reducerPath: "splitApi",
    baseQuery: baseQueryWithRetry,
    refetchOnMountOrArgChange: true,
    endpoints: () => ({}),
})