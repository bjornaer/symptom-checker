import React, { FunctionComponent } from "react";
import { FlowNode } from "typescript";
import logo from './logo.svg';

const Item = (props: {name: string, hpo: string, onClick: React.MouseEventHandler}) => {
  return (
    <div className="item" onClick={props.onClick} id={props.hpo}>
      <div className="logo">
        <img src={logo} alt={props.name} />
      </div>
      <div className="name">
        <p>{props.name}</p>
      </div>
    </div>
  );
};

export default Item;