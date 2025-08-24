import { classNames } from "~/lib/classnames";
import "./Button.css";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "link";
  wide?: boolean;
};

export const Button = ({ children, ...props }: ButtonProps) => {
  const className = classNames(
    props.className,
    "button",
    `button--${props.variant || "primary"}`
  );

  return (
    <button {...props} className={className}>
      {children}
    </button>
  );
};
