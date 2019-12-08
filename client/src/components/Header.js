import React, { useState } from "react";
import { Link, withRouter } from "react-router-dom";
import { Menu } from "semantic-ui-react";
import Cookies from "js-cookie";

const Header = props => {
  const [activeItem, setActiveItem] = useState(props.location.pathname);
  const handleItemClick = (e, { name }) => setActiveItem(name);
  const handleLogout = () => {
    Cookies.remove("auth_token", { path: "/" });
    props.history.push("/");
  };
  return (
    <div>
      <h2 className="ui huge header" style={{ textAlign: "center" }}>
        {"Service Deployment Project Website"}
      </h2>
      <Menu secondary>
        <Menu.Item
          as={Link}
          to="/profile"
          name="/profile"
          active={activeItem.startsWith("/profile")}
          onClick={handleItemClick}
        >
          My Profile
        </Menu.Item>
        <Menu.Item
          as={Link}
          to="/users"
          name="/users"
          active={activeItem.startsWith("/users")}
          onClick={handleItemClick}
        >
          All Students
        </Menu.Item>
        <Menu.Menu position="right">
          <Menu.Item onClick={handleLogout}>Logout</Menu.Item>
        </Menu.Menu>
      </Menu>
    </div>
  );
};

export default withRouter(Header);
