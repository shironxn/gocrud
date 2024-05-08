"use client";

import { AspectRatio } from "@radix-ui/react-aspect-ratio";
import Image from "next/image";

const Header = () => {
  return (
    <div className="w-full">
      <AspectRatio ratio={14 / 2}>
        <Image src="/header.jpg" alt="Image" className="object-cover" fill />
      </AspectRatio>
    </div>
  );
};

export { Header };
