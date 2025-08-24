import "./Message.css";

type MessageProps = {
  children: React.ReactNode;
  type: "error";
};

const ErrIcon = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      className="message__icon"
    >
      <circle cx="12" cy="12" r="10"></circle>
      <line x1="12" x2="12" y1="8" y2="12"></line>
      <line x1="12" x2="12.01" y1="16" y2="16"></line>
    </svg>
  );
};

export const Message = ({ children, type }: MessageProps) => {
  return (
    <div className={`message message--${type}`}>
      {type === "error" && <ErrIcon />}
      <div className="message__content">{children}</div>
    </div>
  );
};

const Title = ({ children }: { children: React.ReactNode }) => {
  return <h3 className="message__title">{children}</h3>;
};

Message.Title = Title;
