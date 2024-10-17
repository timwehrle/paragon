import type { NextPage } from 'next'
import Head from 'next/head'
import { Greet } from '../../wailsjs/wailsjs/go/main/App'
import { ListVolumesJSON } from '../../wailsjs/wailsjs/go/main/App'
import {useEffect, useState} from "react";
import {main} from "../../wailsjs/wailsjs/go/models";
import VolumeInfo = main.VolumeInfo;

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
                    <strong>File System:</strong> {volume.fileSystem}
                </li>
            ))}
        </ul>
    </div>
  )
}

export default Home
