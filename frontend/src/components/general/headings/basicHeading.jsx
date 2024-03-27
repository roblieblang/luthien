export default function BasicHeading({ text, textSize }) {
  return (
    <div
      className={`text-center text-customStroke ${
        textSize === undefined ? "text-4xl" : textSize
      } w-screen mb-10`}
    >
      <h1>{text}</h1>
    </div>
  );
}
