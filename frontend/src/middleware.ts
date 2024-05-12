import { headers } from "next/headers";
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
  }

  if (
    request.cookies.has("refresh-token") &&
    !request.cookies.has("access-token")
  ) {
    const refreshToken = request.cookies.get("refresh-token");
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/refresh`, {
      method: "POST",
      headers: headers(),
    });

    const result = await res.json();
    const accessTokenExpiresInMinutes = 10;
    const accessTokenExpires = new Date(
      Date.now() + accessTokenExpiresInMinutes * 60 * 1000
    );

    console.log(result);

    const response = NextResponse.next();

    response.cookies.set("access-token", result.access_token, {
      httpOnly: true,
      path: "/",
      expires: accessTokenExpires,
    });

    return response;
  }

  if (request.nextUrl.pathname === "/profile") {
    if (!request.cookies.has("refresh-token")) {
      return NextResponse.redirect(new URL("/", request.url));
    }
    return NextResponse.next();
  }

  return NextResponse.next();
}
