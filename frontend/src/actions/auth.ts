"use server";

import { AuthLogin, AuthRegister } from "@/lib/schema/auth";
import { cookies, headers } from "next/headers";

const BASE_API_URL = process.env.NEXT_PUBLIC_API_URL;

const Login = async (data: AuthLogin) => {
  const res = await fetch(`${BASE_API_URL}/auth/login`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json" },
  });

  const result = await res.json();
  if (!res.ok) {
    return result.error;
  }

  const accessTokenExpiresInMinutes = 10;
  const accessTokenExpires = new Date(
    Date.now() + accessTokenExpiresInMinutes * 60 * 1000
  );

  cookies().set("access-token", result.tokens.access_token, {
    httpOnly: true,
    path: "/",
    expires: accessTokenExpires,
  });

  const refreshTokenExpiresInMonths = 1;
  const refreshTokenExpires = new Date(
    Date.now() + refreshTokenExpiresInMonths * 30 * 24 * 60 * 60 * 1000
  );

  cookies().set("refresh-token", result.tokens.refresh_token, {
    httpOnly: true,
    path: "/",
    expires: refreshTokenExpires,
  });
};

const Register = async (data: AuthRegister) => {
  const res = await fetch(`${BASE_API_URL}/auth/register`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json" },
  });

  if (!res.ok) {
    const result = await res.json();
    return result.error;
  }
};

const Logout = async () => {
  const res = await fetch(`${BASE_API_URL}/auth/logout`, {
    method: "POST",
    headers: headers(),
  });

  if (!res.ok) {
    const result = await res.json();
    return result.error;
  }

  cookies().delete("access-token");
  cookies().delete("refresh-token");
};

export { Login, Register, Logout };
