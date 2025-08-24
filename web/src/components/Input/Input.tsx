import React, { type InputHTMLAttributes } from "react";
import { classNames } from "~/lib/classnames";
import "./Input.css";

type InputProps = InputHTMLAttributes<HTMLInputElement>;

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  (props, ref) => {
    const className = classNames(props.className, "input");
    return <input ref={ref} {...props} className={className} />;
  }
);
Input.displayName = "Input";
