// import { useState, useEffect } from "react"

// export default function useFetch<T>(url: string) {

//     const [data, setData] = useState<T|null>(null)
//     const [loading, setLoading] = useState(true)
//     const [error, setError] = useState<string|null>(null)

//     useEffect(() => {
//         async function fetchData() {
//             try {
//                 setLoading(true)
//                 setError(null)

//                 const res = await fetch(url)

//                 if (!res.ok) {
//                     throw new Error("Failed to fetch data");
//                 }

//                 const json = await res.json();
//                 setData(json)

//             } catch(err) {
//                 setError(err instanceof Error ? err.message: "Something wrong")
//             } finally {
//                 setLoading(false)
//             }
//         }
//         fetchData()
//     },[url])

//     return{ data, loading, error };
// }

import { useEffect, useState } from "react";
import { api } from "../lib/api"

export default function useFetch<T>(url:string) {

    const [data, setData] = useState<T|null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string|null>(null);

    useEffect(() => {
        async function fetchData() {
            try {

            setLoading(true)
            setError(null)

            const res = await api.get<T>(url)
            setData(res.data);

            } catch(err) {
                setError(err instanceof Error ? err.message: "Failed something wrong")
            } finally {
                setLoading(false)
            }
        }
        fetchData()
    }, [url])

    return {data, loading, error};
}