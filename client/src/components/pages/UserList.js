import React, { useState, useEffect } from "react";
import axios from "axios";
import { Table } from "semantic-ui-react";

import Header from "../Header";

const UserList = props => {
  const [userData, setUserData] = useState([]);

  useEffect(() => {
    axios.get("/api/v1/priv/users").then(res => {
      setUserData(res.data);
    });
  }, []);

  const renderRow = (el, i) => {
    return (
      <Table.Row key={i}>
        <Table.Cell>{el.first_name}</Table.Cell>
        <Table.Cell>{el.last_name}</Table.Cell>
        <Table.Cell>{el.school_email}</Table.Cell>
      </Table.Row>
    );
  };

  return (
    <div>
      <Header />
      <Table striped>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>First Name</Table.HeaderCell>
            <Table.HeaderCell>Last Name</Table.HeaderCell>
            <Table.HeaderCell>School Email</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>{userData.map((el, i) => renderRow(el, i))}</Table.Body>
      </Table>
    </div>
  );
};

export default UserList;
