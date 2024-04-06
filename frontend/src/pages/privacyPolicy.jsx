/* eslint-disable react/no-unescaped-entities */

export default function PrivacyPolicy() {
  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold">Privacy Policy for Luthien</h1>
      <p className="mt-4">
        <strong>Effective Date:</strong> 3/29/2024
      </p>
      <section className="mt-8">
        <p>
          Welcome to Luthien ("Luthien", "we", "our", or "us"), a mobile
          application designed to convert playlists between Spotify and YouTube.
          Our app requires users to log in to each service, fetches playlists
          from each source service, and uses them to create a new playlist in
          the target service. This Privacy Policy outlines how we collect, use,
          and share information about you through Luthien accessible at{" "}
          <a
            target="_blank"
            href="https://luthien.vercel.app/"
            className="text-blue-600 hover:underline"
            rel="noreferrer"
          >
            https://luthien.vercel.app/
          </a>
          .
        </p>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">
          Information Collection and Use
        </h2>
        <p>
          To provide Luthien's functionality, we collect the following types of
          information:
        </p>
        <ul className="list-none pl-4 mt-2">
          <li>
            <strong>Account Information</strong>
            <ul className="list-none pl-4 mt-1">
              <li>
                Email Address: Used for account registration and communication.
              </li>
              <li>
                Spotify and YouTube Login Data: Used to access your playlists on
                Spotify and YouTube. We do not store your Spotify or YouTube
                passwords.
              </li>
            </ul>
          </li>
          <li className="mt-2">
            <strong>Playlist Data</strong>
            <ul className="list-none pl-4 mt-1">
              <li>
                Playlist Names and Contents: We temporarily fetch the names and
                contents of your playlists to facilitate the conversion process.
              </li>
            </ul>
          </li>
        </ul>
        <p>
          We do not collect detailed usage data such as your device's IP
          address, browser type, or browsing behavior through Luthien.
        </p>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">Purpose of Data Collection</h2>
        <p>The data we collect is used to:</p>
        <ul className="list-none pl-4 mt-2">
          <li>Provide and maintain the functionality of Luthien.</li>
          <li>
            Allow you to access and use Luthien's playlist conversion feature.
          </li>
          <li>Provide customer support and respond to your inquiries.</li>
        </ul>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">Sharing Your Information</h2>
        <p>
          We do not sell, trade, or rent your personal identification
          information to others. Playlist data is processed to perform the
          conversion between services and is not shared with any third parties
          beyond what is necessary to perform the service. We ensure the limited
          use of Google user data, strictly according to the purposes outlined
          in this Privacy Policy.
        </p>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">Security of Data</h2>
        <p>
          The security of your data is important to us. We implement
          commercially acceptable means to protect your personal information,
          though no method of transmission over the Internet or electronic
          storage is 100% secure.
        </p>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">
          Compliance with Google's Limited Use Policy
        </h2>
        <p>
          Luthien's use and transfer of information received from Google APIs
          will adhere to the{" "}
          <a
            href="https://developers.google.com/terms/api-services-user-data-policy"
            target="_blank"
            className="text-blue-600 hover:underline"
            rel="noreferrer"
          >
            Google API Services User Data Policy
          </a>
          , including the Limited Use requirements. By using Luthien, you agree
          to our Privacy Policy, Terms of Service, and are subject to the Google
          API Services User Data Policy. This ensures your data is handled
          securely and in accordance with privacy standards.
        </p>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">
          Changes to This Privacy Policy
        </h2>
        <p>
          We may update our Privacy Policy to reflect changes to our information
          practices. We will notify you of any changes by posting the new
          Privacy Policy on Luthien. We encourage you to periodically review
          this page for the latest information on our privacy practices.
        </p>
      </section>
      <section className="mt-8">
        <h2 className="text-xl font-semibold">Contact Us</h2>
        <p>
          If you have any questions about this Privacy Policy, please contact
          us:
        </p>
        <p>
          By email:{" "}
          <a
            href="mailto:robertlieblang@gmail.com"
            className="text-blue-600 hover:underline"
          >
            robertlieblang@gmail.com
          </a>
        </p>
      </section>
    </div>
  );
}
