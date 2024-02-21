export default function BasicHeading({ text }) {
  return (
    <di className="flex items-center justify-center min-h-screen">
      <div className="px-10 mb-2 text-center text-customStroke text-4xl w-screen">
        <h1>{text}</h1>
      </div>
    </di>
    
  );
}
