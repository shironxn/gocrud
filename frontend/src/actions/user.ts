"use server";

import { headers } from "next/headers";
import { notFound } from "next/navigation";

const BASE_API_URL = process.env.NEXT_PUBLIC_API_URL;

const GetUserMe = async () => {
  try {
    const res = await fetch(`${BASE_API_URL}/users/me`, {
      headers: headers(),
      cache: "no-store",
    });

    const result = await res.json();
    if (!res.ok) {
      return { error: result?.error };
    }

    return { data: result };
  } catch (error) {
    console.error("Error during fetching user data:", error);
    throw error;
  }
};

const GetUserByName = async (name: string) => {
  try {
    const res = await fetch(`${BASE_API_URL}/users?name=${name}&details=true`, {
      cache: "no-store",
    });

    if (res.status === 404) {
      notFound();
    }

    const result = await res.json();
    if (!res.ok) {
      return { error: result?.error };
    }

    return await result.users[0];
  } catch (error) {
    console.error("Error during fetching user by ID:", error);
    throw error;
  }
};

export { GetUserMe, GetUserByName };
