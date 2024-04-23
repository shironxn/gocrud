import { AspectRatio } from "@/components/ui/aspect-ratio";
import Image from "next/image";
import { LoginForm } from "@/components/auth-form";

export default async function Page() {
  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="max-w-6xl flex flex-col space-y-5 p-5 md:flex-row w-full md:space-x-10">
        <div className="w-full md:w-1/2 flex justify-center items-center ">
          <AspectRatio ratio={1 / 1}>
            <Image
              src="/login.png"
              alt="Image"
              className="rounded-md object-cover"
              fill
            />
          </AspectRatio>
        </div>
        <div className="w-full space-y-6 md:w-1/2">
          <div>
            <h1 className="text-center scroll-m-20 text-5xl font-extrabold tracking-tight md:text-5xl">
              Login
            </h1>
          </div>
          <LoginForm />
        </div>
      </div>
    </div>
  );
}
