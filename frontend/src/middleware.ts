import axios from "axios";
import { NextResponse } from "next/server";
import { NextRequest } from "next/server";

export async function middleware(request: NextRequest) {
  if (
    request.nextUrl.pathname.startsWith("/login") ||
    request.nextUrl.pathname.startsWith("/register")
  ) {
    if (request.cookies.has("refresh-token")) {
      return NextResponse.redirect(new URL("/", request.url));
    }

    if (
      request.cookies.has("refresh-token") &&
      !request.cookies.has("access-token")
    ) {
      const refreshToken = request.cookies.get("refresh-token");
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/auth/refresh`,
        {
          body: refreshToken?.value,
          method: "POST",
        }
      );
      const result = await res.json();
      console.log(result);
    }
    return NextResponse.next();
  }

  return NextResponse.next();
}
