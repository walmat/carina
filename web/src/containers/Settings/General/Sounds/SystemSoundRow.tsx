import { Dispatch, SetStateAction, useCallback } from "react";
import { useDispatch } from "react-redux";
import { useDropzone } from "react-dropzone";
import { AnimatePresence, motion } from "framer-motion";
import styled from "styled-components";
import { VolumeX, Volume2, Play, Square } from "react-feather";
import MultiSlider, { Progress, Dot } from "react-multi-bar-slider";

import { Buttons, Typography } from "../../../../elements";
import { editSound, revertSound, changeVolume } from "../../../../stores/Main/reducers/settings";

interface SystemSoundRowProps {
  soundType: string;
  selector: {
    name: string;
    volume: number;
    default: boolean;
  };
  isPlaying: string;
  setIsPlaying: Dispatch<SetStateAction<string>>;
  theme: number;
}

const trimFileName = (name: string) => {
  if (name.length <= 13) {
    return name;
  }

  const parts = name.split(".");
  const { length } = parts[parts.length - 1];

  /* subtract 3 ellipsis plus the length of the ext from 12 */
  const leftover = 13 - (length + 3);
  const first = name.slice(0, leftover / 2);
  const last = name.slice(name.length - length - leftover / 2, name.length);

  return `${first}...${last}`;
};

const playAudio = (type: string, play = true) => {
  switch (type) {
    default:
    case "checkout":
      if (play) {
        return window.RPCAction('preferences:playCheckout')
      }
      return window.RPCAction('preferences:stopAllSounds')
    case "harvester":
      if (play) {
        return window.RPCAction('preferences:playHarvester')
      }
      return window.RPCAction('preferences:stopAllSounds')
  }
};

const uploadAudio = (type: string, path = "") => {
  switch (type) {
    default:
    case "checkout":
      return window.RPCAction('preferences:setCheckoutSound', [path])
    case "harvester":
      return window.RPCAction('preferences:setHarvesterSound', [path])
  }
};

const setVolume = (type: string, volume: number) => {
  switch (type) {
    default:
    case "checkout":
      return window.RPCAction('preferences:setCheckoutVolume', [volume])
    case "harvester":
      return window.RPCAction('preferences:setHarvesterVolume', [volume])
  }
};

const getThemeColors = (theme: number) => {
  if (theme === 1) {
    /* dark mode */
    return {
      progress: "rgba(120, 110, 242, 0.2)",
      backgroundColor: "rgba(255, 255, 255, 0.1)",
    };
  }

  /* light mode */
  return {
    progress: "#EFEDFE",
    backgroundColor: "rgba(0, 0, 0, 0.1)",
  };
};

const SystemSoundRow = ({
  soundType,
  selector,
  theme,
  isPlaying,
  setIsPlaying,
}: SystemSoundRowProps) => {
  const dispatch = useDispatch();
  const colors = getThemeColors(theme);

  const edit = (name: string, type: string) => {
    dispatch(
      editSound({
        name,
        type,
        volume: 100,
      })
    );
  };

  /*
    Change the sound once uploaded. This triggers for both
    drag & drop events and file picker uploads.
  */
  const onDrop = useCallback((acceptedFiles) => {
    acceptedFiles.forEach((file: any) => {
      console.log(file);
      uploadAudio(soundType.toLowerCase(), file.path);
      edit(file.name, soundType.toLowerCase());
    });
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    maxFiles: 1,
    multiple: false,
    accept: [".mp3", ".wav", ".ogg", ".flac"],
    onDrop,
    noClick: true,
    noKeyboard: true,
  });

  /*
    When reverted the sound should be stopped
    and then we can tell Redux to revert the sound.
    We use source.disconnect instead of source.stop because
    the start function can't be called after being stopped.
  */
  const handleRevert = () => {
    if (isPlaying) {
      playAudio(soundType.toLowerCase(), false);
      setIsPlaying("");
    }

    uploadAudio(soundType.toLowerCase());
    dispatch(revertSound({ type: soundType.toLowerCase() }));
  };

  const handleUpload = async () => {
    try {
      const path = await window.RPCAction('fs:selectFileFiltered', [
        "Choose Audio File",
        "Audio File",
        "*.mp3",
        "*.wav",
        "*.ogg",
        "*.flac"
      ]);

      const parts = path.split("/");
      const name = parts[parts.length - 1];

      if (!name) {
        return;
      }

      uploadAudio(soundType.toLowerCase(), path);
      dispatch(editSound({ type: soundType.toLowerCase(), name, volume: 100 }));
    } catch (e) {
      // noop..
    }
  };

  const handleAudio = () => {
    if (isPlaying === soundType.toLowerCase()) {
      playAudio(soundType.toLowerCase(), false);
      setIsPlaying("");
      return;
    }

    setIsPlaying(soundType.toLowerCase());
    playAudio(soundType.toLowerCase());
  };

  const onChangeVolume = (progress: number) => {
    setVolume(soundType.toLowerCase(), progress);
    dispatch(changeVolume({ type: soundType.toLowerCase(), volume: progress }));
  };

  return (
    <Container {...getRootProps()}>
      <AnimatePresence exitBeforeEnter>
        {isDragActive && (
          <Dropzone
            key="uploader"
            initial={{ opacity: 0 }}
            animate={{ opacity: isDragActive ? 1 : 0 }}
            transition={{ duration: 0.35 }}
            exit={{ opacity: 0 }}
          >
            <p>Upload File</p>
          </Dropzone>
        )}
        {!isDragActive && (
          <Col flex={1} ai="stretch" jc="flex-start">
            <Flex
              key="file"
              m="0 0 16px 0"
              initial={{ opacity: 1 }}
              animate={{ opacity: !isDragActive ? 1 : 0 }}
              transition={{ duration: 0.15 }}
              exit={{ opacity: 0 }}
            >
              <Col ai="center" jc="center">
                <Title>{soundType}</Title>
              </Col>
              <Col ai="center" jc="center" onClick={handleAudio}>
                <File isPlaying={isPlaying === soundType.toLowerCase()}>
                  {isPlaying === soundType.toLowerCase() ? (
                    <StopIcon />
                  ) : (
                    <PlayIcon />
                  )}
                  {trimFileName(selector?.name || "")}
                </File>
              </Col>
              <ColFill>
                <RowEnd style={{ alignItems: "center" }}>
                  <RevertButton
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    disabled={selector?.default}
                    type="button"
                    onClick={handleRevert}
                  >
                    Revert
                  </RevertButton>
                  <Buttons.Primary
                    width="9.5vw"
                    height={32}
                    text="Upload File"
                    onClick={handleUpload}
                  />
                  <input {...getInputProps()} />
                </RowEnd>
              </ColFill>
            </Flex>
            <Flex>
              <Col ai="center" jc="center">
                <Muted onClick={() => onChangeVolume(0)} />
              </Col>
              <Col m="0 16px" flex={1} ai="center" jc="center">
                <MultiSlider
                  height={8}
                  slidableZoneSize={10}
                  backgroundColor={colors.backgroundColor}
                  onSlide={onChangeVolume}
                  roundedCorners
                >
                  <Progress
                    color={
                      selector.volume === 0
                        ? "rgba(255, 255, 255, 0.1)"
                        : colors.progress
                    }
                    progress={selector.volume}
                  >
                    <Handle
                      width={14}
                      height={14}
                      color={selector.volume === 0 ? "#616161" : "#786EF2"}
                    />
                  </Progress>
                </MultiSlider>
              </Col>
              <Col ai="center" jc="center">
                <Max onClick={() => onChangeVolume(100)} />
              </Col>
            </Flex>
          </Col>
        )}
      </AnimatePresence>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  align-items: center;
  flex: 1;
`;

const RevertButton = styled(motion.button)<{ disabled: boolean }>`
  font-size: 14px;
  font-weight: 700;
  background: transparent;
  color: ${({ theme }) => theme.colors.failed};
  width: 7.5vw;
  opacity: ${({ disabled }) => (disabled ? 0.4 : 1)};
  cursor: ${({ disabled }) => (disabled ? "default" : "cursor")};
  border: none;
  margin-right: 8px;
  margin-left: auto;
  border-radius: 4px;
`;

const Dropzone = styled(motion.div)`
  border: 1px dashed ${({ theme }) => theme.colors.primary};
  margin: 16px 0;
  border-radius: 4px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 14px;
  flex: 1;
  height: 32px;
  background-color: ${({ theme }) => theme.colors.secondary};
  color: ${({ theme }) => theme.colors.primary};
  cursor: pointer;
`;

const Flex = styled(motion.div)<{ m?: string }>`
  display: flex;
  flex: 1;
  ${({ m }) =>
    m
      ? `
    margin: ${m};
  `
      : ""}
`;

const Title = styled(Typography.H2)`
  font-size: 14px;
  min-width: 62px;
  margin: 0 16px 0 0;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const File = styled(Typography.H2)<{ isPlaying: boolean }>`
  font-size: 14px;
  margin: 0;
  display: flex;
  align-items: flex;
  justify-content: flex;
  border-radius: 22.5px;
  padding: 8px 16px;
  font-weight: 500;
  cursor: pointer;
  opacity: ${({ isPlaying }) => (isPlaying ? 0.75 : 1)};
  background-color: ${({ theme, isPlaying }) =>
    isPlaying ? theme.colors.failed : theme.colors.secondary};
  color: ${({ theme, isPlaying }) =>
    isPlaying ? "white" : theme.colors.primary};
`;

const Col = styled.div<{ flex?: number; m?: string; ai?: string; jc?: string }>`
  display: flex;
  flex-direction: column;
  align-items: ${({ ai }) => ai};
  justify-content: ${({ jc }) => jc};
  ${({ flex }) =>
    flex
      ? `
    flex: ${flex};
  `
      : ""}
  ${({ m }) =>
    m
      ? `
    margin: ${m};
  `
      : ""}
`;

const ColFill = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const RowEnd = styled.div`
  display: flex;
  justify-content: flex-end;
  flex: 1;
`;

const Muted = styled(VolumeX)`
  width: 14px;
  height: 14px;
  color: ${({ theme }) => theme.colors.h1};
`;

const Max = styled(Volume2)`
  width: 14px;
  height: 14px;
  color: ${({ theme }) => theme.colors.h1};
`;

const StopIcon = styled(Square)`
  width: 14px;
  height: 14px;
  margin: auto 4px auto 0;
  fill-opacity: 0.4;
  stroke: white;
  fill: white;
`;

const PlayIcon = styled(Play)`
  width: 14px;
  height: 14px;
  margin: auto 4px auto 0;
  stroke: ${({ theme }) => theme.colors.primary};
  fill: ${({ theme }) => theme.colors.secondary};
`;

const Handle = styled(Dot)`
  cursor: pointer;
`;

export default SystemSoundRow;
