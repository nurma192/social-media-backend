import {api} from "./api";
import type {Post} from "../types";

export const postApi = api.injectEndpoints({
    endpoints: (builder) => ({
        createPost: builder.mutation<Post, {content: string}>({
            query: (postData) => ({
                url: '/posts',
                method: "POST",
                body: postData
            })
        }),
        getAllPosts: builder.query<Post[], void>({
            query: () => ({
                url: '/posts',
                method: "GET"
            })
        }),
        getPostByID: builder.query<Post, string>({
            query: (id) => ({
                url: `/api/posts/${id}`,
                method: "GET"
            })
        }),
        deletePost: builder.mutation<void, string>({
            query: (id) => ({
                url: `/api/posts/${id}`,
                method: "DELETE"
            })
        })
    })
})

export const {
    useCreatePostMutation,
    useGetAllPostsQuery,
    useGetPostByIDQuery,
    useDeletePostMutation,
    useLazyGetAllPostsQuery,
    useLazyGetPostByIDQuery
} = postApi

export const {
    endpoints: { createPost, getAllPosts, getPostByID, deletePost },
} = postApi