import { GetUserMe } from "@/actions/user";
import { NavMenu } from "./nav-menu";

const Navbar = async () => {
  const result = await GetUserMe();

  return <NavMenu user={result} />;
};

export { Navbar };
