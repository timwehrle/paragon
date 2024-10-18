import type { NextPage } from "next";
import { ListVolumesJSON } from "../../wailsjs/wailsjs/go/app/App";
import { useEffect, useState } from "react";
import { app } from "../../wailsjs/wailsjs/go/models";
import VolumeInfo = app.VolumeInfo;

const Home: NextPage = () => {
  const [volumes, setVolumes] = useState<VolumeInfo[]>([]);

  useEffect(() => {
    async function fetchVolumes() {
      try {
        const result = await ListVolumesJSON();
        const parsedResult = JSON.parse(result);
        setVolumes(parsedResult);
      } catch (error) {
        console.error("Error fetching volumes:", error);
      }
    }

    fetchVolumes();
  }, []);
  return (
    <div>
      <h1>Test</h1>
      <ul>
        {volumes.map((volume, index) => (
          <li key={index}>
            <strong>Drive:</strong> {volume.drive} <br />
            <strong>Type:</strong> {volume.type} <br />
            <strong>Volume Label:</strong> {volume.volumeLabel} <br />
            <strong>File System:</strong> {volume.fileSystem} <br />
            <strong>Serial Number:</strong> {volume.serialNumber} <br />
            <strong>Max Component Size:</strong> {volume.maxComponentSize}{" "}
            <br />
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Home;
