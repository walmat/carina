const attachRpc = async () =>
  new Promise((resolve) =>
    document.addEventListener("astilectron-ready", () => {
      if (!(window as any).RPCAction) {
        (window as any).RPCAction = (action: string, input: any = []) => {
          console.log(action, input);
          return new Promise((resolve, reject) => {
            astilectron.sendMessage(
              {
                type: "rpc",
                data: JSON.stringify({
                  action,
                  input: [...input],
                }),
              },
              (response: { output: any; error: string | undefined }) => {
                if (!response) {
                  return;
                }

                // TODO: check why some functions return an empty object when error is nil
                if (response.error && JSON.stringify(response.error) != "{}") {
                  return reject(response.error);
                }

                return resolve(response.output);
              }
            );
          });
        };
      }

      resolve(true);
    })
  );

export { attachRpc };
